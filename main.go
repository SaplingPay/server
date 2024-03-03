package main

import (
	"log"
	"os"
	"strings"

	"github.com/SaplingPay/server/db"
	"github.com/SaplingPay/server/handlers"
	"github.com/SaplingPay/server/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting server")

	// Initialize the Gin engine
	r := gin.Default()

	env := os.Getenv("ENV")
	if env != "production" {
		log.Println("Loading .env file")
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	// Fetch the allowed origins from the environment variable
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsStr == "" {
		log.Println("No allowed origins specified in the environment. Exiting...")
		return
	}
	// Split the environment variable into a slice of allowed origins
	allowedOrigins := strings.Split(allowedOriginsStr, ",")

	// Apply the AllowedOriginsMiddleware globally
	r.Use(middleware.AllowedOriginsMiddleware(allowedOrigins))

	r.Use(cors.Default())

	// Get the MongoDB URI from the .env file
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not found in .env file")
	}

	db.ConnectMongo(mongoURI)

	handlers.SetUpRoutes(r)

	// Start the server
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
