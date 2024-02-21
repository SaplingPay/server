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
	// Initialize the Gin engine
	r := gin.Default()

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the MongoDB URI from the .env file
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not found in .env file")
	}

	db.Connect(mongoURI)

	setUpRoutes(r)

	// Start the server
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}

func setUpRoutes(r *gin.Engine) {

	menuRoutes := r.Group("/menu")
	{
		// Routes for Menu operations
		menuRoutes.GET("/", handlers.GetAllMenus) // Assuming there's a GetAllMenus function
		menuRoutes.POST("/", handlers.CreateMenu)
		menuRoutes.GET("/:menuId", handlers.GetMenu)
		menuRoutes.PUT("/:menuId", handlers.UpdateMenu)
		menuRoutes.DELETE("/:menuId", handlers.DeleteMenu)

		// Nested routes for Menu Item operations under a specific Menu
		menuItemRoutes := menuRoutes.Group("/:menuId/items")
		{
			menuItemRoutes.POST("/", handlers.CreateMenuItem)
			menuItemRoutes.GET("/", handlers.GetAllMenuItems)
			menuItemRoutes.GET("/:itemId", handlers.GetMenuItem)
			menuItemRoutes.PUT("/:itemId", handlers.UpdateMenuItem)
			menuItemRoutes.DELETE("/:itemId", handlers.DeleteMenuItem)
		}
	}

	orderRoutes := r.Group("/orders")
	{
		// Routes for Order operations
		orderRoutes.GET("/", handlers.GetAllOrders)
		orderRoutes.POST("/", handlers.CreateOrder)
		orderRoutes.GET("/:orderId", handlers.GetOrder)
		orderRoutes.PUT("/:orderId", handlers.UpdateOrder)
		orderRoutes.DELETE("/:orderId", handlers.DeleteOrder)

		// Nested routes for Order Item operations under a specific Order
		orderItemRoutes := orderRoutes.Group("/:orderId/items")
		{
			orderItemRoutes.POST("/", handlers.AddOrderItem)
			orderItemRoutes.GET("/", handlers.GetOrderItems)
			orderItemRoutes.GET("/:itemId", handlers.GetOrderItem)
			orderItemRoutes.PUT("/:itemId", handlers.UpdateOrderItem)
			orderItemRoutes.DELETE("/:itemId", handlers.DeleteOrderItem)
		}
	}

	kitchenOrderRoutes := r.Group("/kitchen_orders")
	{
		// Routes for KitchenOrder operations
		kitchenOrderRoutes.GET("/", handlers.GetAllKitchenOrders)
		kitchenOrderRoutes.POST("/", handlers.CreateKitchenOrder)
		kitchenOrderRoutes.GET("/:orderId", handlers.GetKitchenOrder)
		kitchenOrderRoutes.PUT("/:orderId", handlers.UpdateKitchenOrder)
		kitchenOrderRoutes.DELETE("/:orderId", handlers.DeleteKitchenOrder)
	}
}
