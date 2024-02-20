package main

import (
	"github.com/SaplingPay/server/db"
	"github.com/SaplingPay/server/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the Gin engine
	r := gin.Default()

	// Connect to MongoDB
	db.Connect("mongodb://localhost:27017")

	// Setup routes
	r.POST("/menu", handlers.AddMenu)
	r.GET("/menu", handlers.GetMenus)

	// Start the server
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
