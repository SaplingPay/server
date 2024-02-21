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

// CreateKitchenOrder creates a new kitchen order in the database
func CreateKitchenOrder(c *gin.Context) {
	var kitchenOrder models.KitchenOrder

	if err := c.ShouldBindJSON(&kitchenOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if kitchenOrder.Items == nil {
		kitchenOrder.Items = []models.OrderItem{}
	}

	result, err := db.DB.Collection("kitchen_orders").InsertOne(ctx, kitchenOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateKitchenOrder updates an existing kitchen order in the database
func UpdateKitchenOrder(c *gin.Context) {
	orderId := c.Param("orderId") // Get the ID from the URL parameter

	var kitchenOrder models.KitchenOrder
	if err := c.ShouldBindJSON(&kitchenOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert the string ID to MongoDB's ObjectID
	objID, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	// Replace the document based on the ObjectID
	result, err := db.DB.Collection("kitchen_orders").ReplaceOne(ctx, bson.M{"_id": objID}, kitchenOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if a document was actually modified
	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen order not found"})
		return
	}

	c.JSON(http.StatusOK, kitchenOrder)
}

// DeleteKitchenOrder deletes a kitchen order from the database
func DeleteKitchenOrder(c *gin.Context) {
	// Fetching the order ID from the URL parameter
	orderId := c.Param("orderId")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Assuming `orderId` needs to be converted to an ObjectID if you're using MongoDB's default ObjectID
	// If your ID is a string in the database, you can directly use orderId in the filter
	objID, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	_, err = db.DB.Collection("kitchen_orders").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kitchen order deleted successfully"})
}

// GetKitchenOrder retrieves a single kitchen order from the database
func GetKitchenOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert the ID from the URL parameter to an ObjectID
	orderId := c.Param("orderId")
	objID, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var kitchenOrder models.KitchenOrder
	// Use the ObjectID to find the kitchen order
	if err := db.DB.Collection("kitchen_orders").FindOne(ctx, bson.M{"_id": objID}).Decode(&kitchenOrder); err != nil {
		// Adjust the error handling to distinguish not found errors from other errors
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, kitchenOrder)
}

// GetAllKitchenOrders retrieves all kitchen orders from the database
func GetAllKitchenOrders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var kitchenOrders []models.KitchenOrder
	cursor, err := db.DB.Collection("kitchen_orders").Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var kitchenOrder models.KitchenOrder
		if err := cursor.Decode(&kitchenOrder); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		kitchenOrders = append(kitchenOrders, kitchenOrder)
	}

	c.JSON(http.StatusOK, kitchenOrders)
}
