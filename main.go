package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// generate hello world for main.go
func main() {
	fmt.Println("Hello, World!")
	router := gin.Default()

	// Define a GET route
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// Start the server on port 8080
	router.Run(":8080")
}
