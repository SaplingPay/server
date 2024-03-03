package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/SaplingPay/server/db"
	"github.com/SaplingPay/server/handlers"
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
	r.Use(AllowedOriginsMiddleware(allowedOrigins))

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

func AllowedOriginsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		referer := c.Request.Header.Get("Referer")
		clientIP := c.ClientIP()

		log.Println("Request Origin:", origin)
		log.Println("Request Referer:", referer)
		log.Println("Client IP:", clientIP)

		var allowed bool

		if origin != "" && matchAllowedOrigin(origin, allowedOrigins) {
			log.Println("Origin allowed:", origin)
			allowed = true
		} else if referer != "" && matchAllowedOrigin(referer, allowedOrigins) {
			log.Println("Referer allowed:", referer)
			allowed = true
		} else if isLocalhost(clientIP) {
			log.Println("Localhost allowed:", clientIP)
			allowed = true
		}

		if allowed {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Origin not allowed"})
		}
	}
}

func matchAllowedOrigin(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}

func isLocalhost(ip string) bool {
	return ip == "127.0.0.1" || ip == "::1"
}
