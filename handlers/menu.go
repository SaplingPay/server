package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/SaplingPay/server/db"
	"github.com/SaplingPay/server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateMenu creates a new menu in the database
func CreateMenu(c *gin.Context) {
	var menu models.Menu

	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if menu.Items == nil {
		menu.Items = []models.MenuItem{}
	}

	result, err := db.DB.Collection("menus").InsertOne(ctx, menu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateMenu updates an existing menu in the database
func UpdateMenu(c *gin.Context) {
	menuID := c.Param("menuId") // Get the ID from the URL parameter

	var menu models.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert the string ID to MongoDB's ObjectID
	objID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	// Replace the document based on the ObjectID
	result, err := db.DB.Collection("menus").ReplaceOne(ctx, bson.M{"_id": objID}, menu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if a document was actually modified
	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// DeleteMenu deletes a menu from the database
func DeleteMenu(c *gin.Context) {
	// Fetching the menu ID from the URL parameter
	menuID := c.Param("menuId")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Assuming `menuID` needs to be converted to an ObjectID if you're using MongoDB's default ObjectID
	// If your ID is a string in the database, you can directly use menuID in the filter
	objID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	_, err = db.DB.Collection("menus").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu deleted successfully"})
}

// GetMenu retrieves a single menu from the database
func GetMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert the ID from the URL parameter to an ObjectID
	menuID := c.Param("menuId")
	objID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var menu models.Menu
	// Use the ObjectID to find the menu
	if err := db.DB.Collection("menus").FindOne(ctx, bson.M{"_id": objID}).Decode(&menu); err != nil {
		// Adjust the error handling to distinguish not found errors from other errors
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, menu)
}

// GetAllMenus retrieves all menus from the database
func GetAllMenus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var menus []models.Menu
	cursor, err := db.DB.Collection("menus").Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var menu models.Menu
		if err := cursor.Decode(&menu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		menus = append(menus, menu)
	}

	c.JSON(http.StatusOK, menus)
}
