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
	//r.Use(middleware.AuthMiddleware())

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

	// menuRoutesV2 := r.Group("/menusV2")
	// {
	// 	menuRoutesV2.GET("/", GetAllMenusV2)
	// 	menuRoutesV2.POST("/:venueId", CreateMenuV2)
	// 	menuRoutesV2.GET("/:menuId", GetMenuV2)
	// 	menuRoutesV2.PUT("/:menuId", UpdateMenuV2)
	// 	menuRoutesV2.DELETE("/:menuId", DeleteMenuV2)

	// menuItemRoutesV2 := menuRoutesV2.Group("/:menuId/items")
	// {
	// 	menuItemRoutesV2.POST("/", CreateMenuItemV2)
	// 	menuItemRoutesV2.GET("/", GetAllMenuItemsV2)
	// 	menuItemRoutesV2.GET("/:itemId", GetMenuItemV2)
	// 	menuItemRoutesV2.PUT("/:itemId", UpdateMenuItemV2)
	// 	menuItemRoutesV2.DELETE("/:itemId", DeleteMenuItemV2)
	// }
	// }

	venueRoutes := r.Group("/venues")
	{
		venueRoutes.GET("/", GetAllVenues)
		venueRoutes.POST("/", CreateVenue)
		venueRoutes.GET("/:venueId", GetVenue)
		venueRoutes.PUT("/:venueId", UpdateVenue)
		venueRoutes.DELETE("/:venueId", DeleteVenue)

		venueMenuRoutes := venueRoutes.Group("/:venueId/menu")
		{
			venueMenuRoutes.POST("/", CreateMenuV2)
			venueMenuRoutes.GET("/:menuId", GetMenuV2)
			venueMenuRoutes.PUT("/:menuId", UpdateMenuV2)
			venueMenuRoutes.DELETE("/:menuId", DeleteMenuV2)
		}

		venueMenuItemRoutes := venueRoutes.Group("/:venueId/menu/:menuId/items")
		{
			venueMenuItemRoutes.POST("/", CreateMenuItemV2)
			venueMenuItemRoutes.GET("/", GetAllMenuItemsV2)
			venueMenuItemRoutes.GET("/:itemId", GetMenuItemV2)
			venueMenuItemRoutes.PUT("/:itemId", UpdateMenuItemV2)
			venueMenuItemRoutes.DELETE("/:itemId", DeleteMenuItemV2)
		}
	}

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/", CreateUser)
		userRoutes.GET("/", GetAllUsers)
		userRoutes.GET("/:userId", GetUser)
		userRoutes.PUT("/:userId", UpdateUser)
		userRoutes.DELETE("/:userId", DeleteUser)
	}

	userV2Routes := r.Group("/usersV2")
	{
		userV2Routes.POST("/", CreateUserV2)
		userV2Routes.GET("/", GetAllUsersV2)
		userV2Routes.GET("/:userId", GetUserV2)
		userV2Routes.PUT("/:userId", UpdateUserV2)
		userV2Routes.DELETE("/:userId", DeleteUserV2)
	}

	menuParserRoutes := r.Group("/menuParser")
	{
		menuParserRoutes.POST("/", ParseMenuCard)
	}

	r.GET("/GetMenusByUserID/:userId", GetMenuByUserID)
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
