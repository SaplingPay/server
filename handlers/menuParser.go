package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/SaplingPay/server/models"
	"github.com/SaplingPay/server/repositories"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const loggerTag = "[menu-parser]"

const instructions = `
		You are a parse for restaurant menus. Your job is to return the categories, items, descriptions and prices.
		The menu is in a PDF format, and is attached to this request.
		The returned message should be just a JSON object, in a valid text/json format.
		There's no need for any text, explanation, markdown, or anything else besides the JSON.
		The return format should be JSON, and should be in the following format:
		[
			{
				  "name": "Veggie Pizza",
				  "price": 15.99,
				  "categories": ["Vegetarian", "Pizza", "Main Course"]
			},
			{
				  "name": "Pepperoni Pizza",
				  "price": 18.99,
				  "categories": ["Pizza", "Main Course"]
			}
		]
	`

func uploadFile(client *openai.Client, r io.Reader) (string, error) {
	log.Println(loggerTag, "Uploading file")
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

	file, err := client.CreateFileBytes(context.Background(), openai.FileBytesRequest{
		Name:    "menu.pdf",
		Bytes:   buf.Bytes(),
		Purpose: openai.PurposeAssistants,
	})

	return file.ID, err
}

func createAssistant(client *openai.Client, fileId string) (string, error) {
	name := "menu-parser"
	model := openai.GPT4TurboPreview
	instructions := instructions
	log.Println(loggerTag, "Creating assistant")
	assistant, err := client.CreateAssistant(
		context.Background(),
		openai.AssistantRequest{
			Model: model,
			Name:  &name,
			Tools: []openai.AssistantTool{
				{
					Type: openai.AssistantToolTypeRetrieval,
				},
			},
			FileIDs:      []string{fileId},
			Instructions: &instructions,
		},
	)

	return assistant.ID, err
}

func createThread(client *openai.Client) (string, error) {
	log.Println(loggerTag, "Creating thread")
	thread, err := client.CreateThread(
		context.Background(),
		openai.ThreadRequest{
			Messages: []openai.ThreadMessage{
				{
					Role:    openai.ThreadMessageRoleUser,
					Content: "I would like to parse this menu, that is attached to this request.",
				},
			},
		},
	)

	return thread.ID, err
}

func runThread(client *openai.Client, threadId string, assistantId string) (openai.RunStatus, error) {
	log.Println(loggerTag, "Running thread")
	run, err := client.CreateRun(context.Background(), threadId, openai.RunRequest{
		AssistantID: assistantId,
	})

	for run.Status == openai.RunStatusQueued || run.Status == openai.RunStatusInProgress || run.Status == openai.RunStatusCancelling {
		time.Sleep(1 * time.Second)
		log.Println(loggerTag, "Waiting for thread to complete. Status:", run.Status)
		run, err = client.RetrieveRun(context.Background(), threadId, run.ID)
	}

	log.Println(loggerTag, "Thread completed with status:", run.Status)

	return run.Status, err
}

func getResult(client *openai.Client, threadId string) (string, error) {
	log.Println(loggerTag, "Retrieving result")
	msgs, err := client.ListMessage(context.Background(), threadId, nil, nil, nil, nil)
	if err != nil {
		return "", err
	}

	log.Println(loggerTag, "Message retrieved", msgs)
	msg := msgs.Messages[0]

	return msg.Content[0].Text.Value, nil
}

func cleanup(client *openai.Client, assistantId string, fileId string, threadId string) {
	log.Println(loggerTag, "Cleaning up")
	if threadId != "" {
		client.DeleteThread(context.Background(), threadId)
	}

	if assistantId != "" {
		client.DeleteAssistant(context.Background(), assistantId)
	}

	if fileId != "" {
		client.DeleteFile(context.Background(), fileId)
	}
}

func parseFileUsingGPTAssistant(r io.Reader) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	fileId, err := uploadFile(client, r)

	if err != nil {
		cleanup(client, "", fileId, "")
		return "", err
	}

	assistantId, err := createAssistant(client, fileId)

	if err != nil {
		cleanup(client, assistantId, fileId, "")
		return "", err
	}

	threadId, err := createThread(client)

	if err != nil {
		cleanup(client, assistantId, fileId, threadId)
		return "", err
	}

	status, err := runThread(client, threadId, assistantId)

	if err != nil {
		cleanup(client, assistantId, fileId, threadId)
		return "", err
	}

	var result string
	if status == openai.RunStatusCompleted {
		r, _ := getResult(client, threadId)
		result = r
	}

	cleanup(client, assistantId, fileId, threadId)
	return result, nil
}

func parseImageUsingGPT4Vision(r io.Reader) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

	imgBase64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

	result, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4VisionPreview,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleUser,
				MultiContent: []openai.ChatMessagePart{
					{
						Type: openai.ChatMessagePartTypeText,
						Text: instructions,
					},
					{
						Type: openai.ChatMessagePartTypeImageURL,
						ImageURL: &openai.ChatMessageImageURL{
							URL:    "data:image/jpeg;base64," + imgBase64Str,
							Detail: openai.ImageURLDetailHigh,
						},
					},
				},
			},
		},
		MaxTokens: 3000,
	})

	for _, choice := range result.Choices {
		log.Println(loggerTag, "Choice:", choice.Message.Content)

		for _, part := range choice.Message.MultiContent {
			if part.Type == openai.ChatMessagePartTypeText {
				log.Println(loggerTag, "Partial text:", part.Text)
			}
		}
	}

	return result.Choices[0].Message.Content, err
}

func ParseMenuCard(c *gin.Context) {
	file, _ := c.FormFile("menu")
	openFile, _ := file.Open()

	menuIdStr := c.Param("menuId")

	log.Println(loggerTag, file.Header)
	log.Println(loggerTag, file.Filename)

	contentType := file.Header.Get("Content-Type")
	if contentType != "application/pdf" && contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	var result string
	if contentType == "application/pdf" {
		r, _ := parseFileUsingGPTAssistant(openFile)
		result = r
	} else {
		r, _ := parseImageUsingGPT4Vision(openFile)
		result = r
	}

	result = strings.ReplaceAll(result, "```json", "")
	result = strings.ReplaceAll(result, "```", "")
	resultBytes := []byte(result)
	openFile.Close()

	var menuItems []models.MenuItemV2
	if err := json.Unmarshal(resultBytes, &menuItems); err != nil {
		log.Println(loggerTag, err)
		c.JSON(http.StatusInternalServerError, err)
	}

	menuId, err := primitive.ObjectIDFromHex(menuIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	items, err := repositories.AddAllMenuItems(menuId, menuItems)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, items)
}
