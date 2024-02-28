package handlers

import "github.com/gin-gonic/gin"

func SetUpRoutes(r *gin.Engine) {

	menuRoutes := r.Group("/menus")
	{
		// Routes for Menu operations
		menuRoutes.GET("/", GetAllMenus) // Assuming there's a GetAllMenus function
		menuRoutes.POST("/", CreateMenu)
		menuRoutes.GET("/:menuId", GetMenu)
		menuRoutes.PUT("/:menuId", UpdateMenu)
		menuRoutes.DELETE("/:menuId", DeleteMenu)

		// Nested routes for Menu Item operations under a specific Menu
		menuItemRoutes := menuRoutes.Group("/:menuId/items")
		{
			menuItemRoutes.POST("/", CreateMenuItem)
			menuItemRoutes.GET("/", GetAllMenuItems)
			menuItemRoutes.GET("/:itemId", GetMenuItem)
			menuItemRoutes.PUT("/:itemId", UpdateMenuItem)
			menuItemRoutes.DELETE("/:itemId", DeleteMenuItem)
		}
	}

	orderRoutes := r.Group("/orders")
	{
		// Routes for Order operations
		orderRoutes.GET("/", GetAllOrders)
		orderRoutes.POST("/", CreateOrder)
		orderRoutes.GET("/:orderId", GetOrder)
		orderRoutes.PUT("/:orderId", UpdateOrder)
		orderRoutes.DELETE("/:orderId", DeleteOrder)

		// Nested routes for Order Item operations under a specific Order
		orderItemRoutes := orderRoutes.Group("/:orderId/items")
		{
			orderItemRoutes.POST("/", AddOrderItem)
			orderItemRoutes.GET("/", GetOrderItems)
			orderItemRoutes.GET("/:itemId", GetOrderItem)
			orderItemRoutes.PUT("/:itemId", UpdateOrderItem)
			orderItemRoutes.DELETE("/:itemId", DeleteOrderItem)
		}
	}

	kitchenOrderRoutes := r.Group("/kitchen_orders")
	{
		// Routes for KitchenOrder operations
		kitchenOrderRoutes.GET("/", GetAllKitchenOrders)
		kitchenOrderRoutes.POST("/", CreateKitchenOrder)
		kitchenOrderRoutes.GET("/:orderId", GetKitchenOrder)
		kitchenOrderRoutes.PUT("/:orderId", UpdateKitchenOrder)
		kitchenOrderRoutes.DELETE("/:orderId", DeleteKitchenOrder)
	}

	paymentRoutes := r.Group("/payments")
	{
		// Routes for payment operations
		paymentRoutes.GET("/", GetAllPayments)
		paymentRoutes.POST("/", CreatePayment)
		paymentRoutes.GET("/:paymentId", GetPayment)
		paymentRoutes.PUT("/:paymentId", UpdatePayment)
		paymentRoutes.DELETE("/:paymentId", DeletePayment)
	}

	// Routes for User operations
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
