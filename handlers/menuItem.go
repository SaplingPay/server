package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/SaplingPay/server/db"
	"github.com/SaplingPay/server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateMenuItem(c *gin.Context) {
	log.Println("CreateMenuItem")

	menuID := c.Param("menuId") // Assuming the menu ID is passed as a URL parameter

	var menuItem models.MenuItem
	if err := c.ShouldBindJSON(&menuItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert menuID from string to primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu ID format"})
		return
	}

	// Generate a new ObjectID for the menu item
	menuItem.ID = primitive.NewObjectID()

	// Add the new menu item to the menu document
	filter := bson.M{"_id": objID}
	update := bson.M{"$push": bson.M{"items": menuItem}}
	_, err = db.DB.Collection("menus").UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menuItem)
}

func GetMenuItem(c *gin.Context) {
	log.Println("GetMenuItem")
	menuID := c.Param("menuId")
	menuItemID := c.Param("itemId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert menuID and menuItemID from string to primitive.ObjectID
	objMenuID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu ID format"})
		return
	}
	objMenuItemID, err := primitive.ObjectIDFromHex(menuItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item ID format"})
		return
	}

	var menu models.Menu
	err = db.DB.Collection("menus").FindOne(ctx, bson.M{"_id": objMenuID}).Decode(&menu)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Iterate over the items to find the specific menu item
	for _, item := range menu.Items {
		if item.ID == objMenuItemID {
			c.JSON(http.StatusOK, item)
			return
		}
	}

	// If the item is not found in the loop
	c.JSON(http.StatusNotFound, gin.H{"error": "menu item not found"})
}

func UpdateMenuItem(c *gin.Context) {
	menuID := c.Param("menuId")
	menuItemID := c.Param("itemId")

	var menuItem models.MenuItem
	if err := c.ShouldBindJSON(&menuItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert menuID from string to primitive.ObjectID
	objMenuID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu ID format"})
		return
	}

	// Convert menuItemID from string to primitive.ObjectID for matching in array
	objMenuItemID, err := primitive.ObjectIDFromHex(menuItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item ID format"})
		return
	}

	menuItem.ID = primitive.NilObjectID // Ensure the ID is not included in the update

	// Update the specified menu item within the menu document
	filter := bson.M{"_id": objMenuID, "items._id": objMenuItemID}
	update := bson.M{"$set": bson.M{"items.$": menuItem}}
	_, err = db.DB.Collection("menus").UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menuItem)
}

func DeleteMenuItem(c *gin.Context) {
	menuID := c.Param("menuId")
	menuItemID := c.Param("itemId")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert menuID from string to primitive.ObjectID
	objMenuID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu ID format"})
		return
	}

	// Convert menuItemID from string to primitive.ObjectID for matching in array
	objMenuItemID, err := primitive.ObjectIDFromHex(menuItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item ID format"})
		return
	}

	// Remove the specified menu item from the menu document
	filter := bson.M{"_id": objMenuID}
	update := bson.M{"$pull": bson.M{"items": bson.M{"_id": objMenuItemID}}}
	_, err = db.DB.Collection("menus").UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu item deleted successfully"})
}

func GetAllMenuItems(c *gin.Context) {
	menuID := c.Param("menuId") // Assuming the menu ID is passed as a URL parameter

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var menu models.Menu
	objID, err := primitive.ObjectIDFromHex(menuID) // Convert menuID to ObjectID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu ID format"})
		return
	}

	err = db.DB.Collection("menus").FindOne(ctx, bson.M{"_id": objID}).Decode(&menu)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if menu.Items == nil {
		menu.Items = []models.MenuItem{} // Ensure the response is an empty array rather than null if no items exist
	}

	c.JSON(http.StatusOK, menu.Items)
}
