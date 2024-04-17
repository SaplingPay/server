package main

import (
	"log"
	"os"

	"github.com/SaplingPay/server/db"
	"github.com/SaplingPay/server/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting server")

	r := gin.Default()

	env := os.Getenv("SERVER_ENV")

	if env == "production" {
		log.Println("Running in production")
		gin.SetMode(gin.ReleaseMode)

	} else if env == "development" {
		log.Println("Running in development")
	} else {
		log.Println("No environment set, defaulting to local")
		log.Println("Loading .env file")
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not found in .env file")
	}

	db.ConnectMongo(mongoURI)

	handlers.SetUpRoutes(r)

	// Start the server
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
