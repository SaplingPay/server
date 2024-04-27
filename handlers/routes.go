package handlers

import (
	"log"
	"os"
	"strings"

	"github.com/SaplingPay/server/payments"

	"github.com/SaplingPay/server/middleware"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine) {
	// Set up the routes that require authentication
	//SetUpAuthRoutes(r)

	r.POST("/getToken", middleware.GetToken)

	// Wrap the routes that require authentication in the AuthMiddleware
	r.Use(middleware.AuthMiddleware())

	payments.AddStripRoutes(r)

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

	venueRoutes := r.Group("/venues")
	{
		venueRoutes.GET("/", GetAllVenues)
		venueRoutes.POST("/", CreateVenue)
		venueRoutes.GET("/:venueId", GetVenue)
		venueRoutes.PUT("/:venueId", UpdateVenue)
		venueRoutes.DELETE("/:venueId", SoftDeleteVenue)

		venueMenuRoutes := venueRoutes.Group("/:venueId/menu")
		{
			venueMenuRoutes.POST("/", CreateMenuV2)
			venueMenuRoutes.POST("/parse/", ParseMenuCard)
			venueMenuRoutes.GET("/:menuId", GetMenuV2)
			venueMenuRoutes.PUT("/:menuId", UpdateMenuV2)
			venueMenuRoutes.DELETE("/:menuId", SoftDeleteMenuV2)
		}
		// get all menus for a venue
		venueMenusRoutes := venueRoutes.Group("/:venueId/menus")
		{
			venueMenusRoutes.GET("/", GetMenusByVenueID)
		}

		venueMenuItemRoutes := venueRoutes.Group("/:venueId/menu/:menuId/items")
		{
			venueMenuItemRoutes.POST("/", CreateMenuItemV2)
			venueMenuItemRoutes.GET("/", GetAllMenuItemsV2)
			venueMenuItemRoutes.GET("/:itemId", GetMenuItemV2)
			venueMenuItemRoutes.PUT("/:itemId", UpdateMenuItemV2)
			venueMenuItemRoutes.DELETE("/:itemId", SoftDeleteMenuItemV2)
		}
	}

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/", CreateUser)
		userRoutes.GET("/", GetAllUsers)
		userRoutes.GET("/:userId", GetUser)
		userRoutes.PUT("/:userId", UpdateUser)
		userRoutes.DELETE("/:userId", DeleteUser)
		userRoutes.GET("/:userId/saves", GetUserSaves)
	}

	userV2Routes := r.Group("/usersV2")
	{
		userV2Routes.POST("/", CreateUserV2)
		userV2Routes.GET("/", GetAllUsersV2)
		userV2Routes.GET("/:userId", GetUserV2)
		userV2Routes.PUT("/:userId", UpdateUserV2)
		userV2Routes.DELETE("/:userId", SoftDeleteUserV2)
		userV2Routes.PUT("/:userId/follow/:followingId", FollowUser)
		userV2Routes.PUT("/:userId/unfollow/:followingId", UnFollowUser)
	}

	r.GET("/GetMenusByUserID/:userId", GetMenuByUserID)

	orders := r.Group("/orders")
	{
		orders.POST("/", CreateOrder)
		orders.GET("/:id", GetOrder)
		orders.PUT("/:id", UpdateOrder)
		orders.DELETE("/:id", SoftDeleteOrder)
		orders.GET("/", GetAllOrders)
	}

	payments := r.Group("/payments")
	{
		payments.POST("/", CreatePayment)
		payments.GET("/:id", GetPayment)
		payments.PUT("/:id", UpdatePayment)
		payments.DELETE("/:id", SoftDeletePayment)
		payments.GET("/", GetAllPayments)
	}
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
