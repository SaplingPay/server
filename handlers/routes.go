package handlers

import (
	"log"
	"os"
	"strings"

	"github.com/SaplingPay/server/middleware"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine) {
	// Set up the routes that require authentication
	// SetUpAuthRoutes(r)

	r.POST("/getToken", middleware.GetToken)

	// Wrap the routes that require authentication in the AuthMiddleware
	r.Use(middleware.AuthMiddleware())

	menuRoutes := r.Group("/menus")
	{
		menuRoutes.GET("/", GetAllMenus)
		menuRoutes.POST("/", CreateMenu)
		menuRoutes.GET("/:menuId", GetMenu)
		menuRoutes.PUT("/:menuId", UpdateMenu)
		menuRoutes.DELETE("/:menuId", DeleteMenu)
		menuRoutes.PUT("/archive/:menuId", ArchiveMenu)

		menuItemRoutes := menuRoutes.Group("/:menuId/items")
		{
			menuItemRoutes.POST("/", CreateMenuItem)
			menuItemRoutes.GET("/", GetAllMenuItems)
			menuItemRoutes.GET("/:itemId", GetMenuItem)
			menuItemRoutes.PUT("/:itemId", UpdateMenuItem)
			menuItemRoutes.DELETE("/:itemId", DeleteMenuItem)
			menuItemRoutes.PUT("/archive/:itemId", ArchiveMenuItem)
		}
	}

	orderRoutes := r.Group("/orders")
	{
		orderRoutes.GET("/", GetAllOrders)
		orderRoutes.POST("/", CreateOrder)
		orderRoutes.GET("/:orderId", GetOrder)
		orderRoutes.PUT("/:orderId", UpdateOrder)
		orderRoutes.DELETE("/:orderId", DeleteOrder)

		orderItemRoutes := orderRoutes.Group("/:orderId/items")
		{
			orderItemRoutes.POST("/", AddOrderItem)
			orderItemRoutes.GET("/", GetOrderItems)
			orderItemRoutes.GET("/:itemId", GetOrderItem)
			orderItemRoutes.PUT("/:itemId", UpdateOrderItem)
			orderItemRoutes.DELETE("/:itemId", DeleteOrderItem)
		}
	}

	paymentRoutes := r.Group("/payments")
	{
		paymentRoutes.GET("/", GetAllPayments)
		paymentRoutes.POST("/", CreatePayment)
		paymentRoutes.GET("/:paymentId", GetPayment)
		paymentRoutes.PUT("/:paymentId", UpdatePayment)
		paymentRoutes.DELETE("/:paymentId", DeletePayment)
	}

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/", CreateUser)
		userRoutes.GET("/", GetAllUsers)
		userRoutes.GET("/:userId", GetUser)
		userRoutes.PUT("/:userId", UpdateUser)
		userRoutes.DELETE("/:userId", DeleteUser)
	}

	r.GET("/GetMenusByUserID/:userId", GetMenuByUserID)

	// Route for uploading files to Supabase storage
	// r.POST("/upload", UploadToSupabase)
}

// func setup auth group
func SetUpAuthRoutes(r *gin.Engine) {
	basicAuthAccount := os.Getenv("AUTH_ACCOUNTS")
	if basicAuthAccount == "" {
		log.Fatal("Basic Auth Account not found in .env file")
	}
	parts := strings.Split(basicAuthAccount, ",")

	authGroup := r.Group("/",
		gin.BasicAuth(gin.Accounts{
			parts[0]: parts[1],
		}))

	authGroup.POST("/getToken", middleware.GetToken)
}
