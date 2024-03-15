package handlers

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const loggerTag = "[menu-parser]"

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
	instructions := `
		You are a parse for restaurant menus. Your job is to return the categories, items, descriptions and prices.
		The menu is in a PDF format, and is attached to this request.
		The returned message should be just a JSON object, in a valid text/json format.
		There's no need for any text, explanation, markdown, or anything else besides the JSON.
		The return format should be JSON, and should be in the following format:
		{	
			"categories": [
				{	
					"name": "Appetizers",	
					"items": [	
						{		
							"name": "Spring Rolls",	
							"description": "Crispy spring rolls filled with vegetables",					
							"price": 5.99
						}
					]
				}
			]
		}
	`
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

func parseFileUsingGPT(r io.Reader) (string, error) {
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

func ParseMenuCard(c *gin.Context) {
	file, _ := c.FormFile("menu")
	openFile, _ := file.Open()

	log.Println(file.Header)
	log.Println(file.Filename)

	contentType := file.Header.Get("Content-Type")
	if contentType != "application/pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	result, _ := parseFileUsingGPT(openFile)
	// TODO - implement better json cleanup
	result = strings.ReplaceAll(result, "```json", "")
	result = strings.ReplaceAll(result, "```", "")
	resultBytes := []byte(result)
	openFile.Close()

	c.Data(http.StatusOK, gin.MIMEJSON, resultBytes)
}
