package handlers

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/SaplingPay/server/db"
	"github.com/gin-gonic/gin"
)

func UploadToSupabase(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	log.Println("File name:", header.Filename)
	// Create a temporary file to store the uploaded file
	tempFile, err := os.CreateTemp("temp-files", header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tempFile.Close()

	log.Println("Temp file:", tempFile.Name())
	// Copy the file data to the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		log.Println("Error copying file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("something")
	// Upload the file to Supabase storage
	uploadResult, err := db.Supabase.UploadFile("menu-assets", tempFile.Name(), tempFile)
	if err != nil {
		log.Println("Error uploading file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Upload result:", uploadResult)
	supabaseUrl := os.Getenv("SUPABASE_URL")
	// Construct the URL of the uploaded file
	fileURL := supabaseUrl + "/storage/v1/object/public/" + uploadResult.Key

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"url":     fileURL,
	})
}
