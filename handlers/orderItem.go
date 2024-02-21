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

func AddOrderItem(c *gin.Context) {
	orderId := c.Param("orderId")
	var orderItem models.OrderItem

	if err := c.ShouldBindJSON(&orderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objOrderId, _ := primitive.ObjectIDFromHex(orderId)
	update := bson.M{"$push": bson.M{"items": orderItem}}
	_, err := db.DB.Collection("orders").UpdateOne(ctx, bson.M{"_id": objOrderId}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderItem)
}

func UpdateOrderItem(c *gin.Context) {
	orderId := c.Param("orderId")
	itemIndex := c.Param("itemIndex") // Assume you're using the array index or change to itemId if using unique IDs for items.
	var orderItem models.OrderItem

	if err := c.ShouldBindJSON(&orderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objOrderId, _ := primitive.ObjectIDFromHex(orderId)
	// Using positional operator $ to update the specific item in the array. This requires MongoDB 3.6+
	update := bson.M{"$set": bson.M{"items." + itemIndex: orderItem}}
	_, err := db.DB.Collection("orders").UpdateOne(ctx, bson.M{"_id": objOrderId}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderItem)
}

func DeleteOrderItem(c *gin.Context) {
	orderId := c.Param("orderId")
	itemId := c.Param("itemId")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objOrderId, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID format"})
		return
	}

	objItemId, err := primitive.ObjectIDFromHex(itemId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID format"})
		return
	}

	filter := bson.M{"_id": objOrderId}
	update := bson.M{"$pull": bson.M{"items": bson.M{"_id": objItemId}}}
	_, err = db.DB.Collection("orders").UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order item deleted successfully"})
}

func GetOrderItems(c *gin.Context) {
	orderId := c.Param("orderId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objOrderId, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID format"})
		return
	}

	var order models.Order
	err = db.DB.Collection("orders").FindOne(ctx, bson.M{"_id": objOrderId}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, order.Items)
}

// GetOrderItem retrieves a single order item from the database
func GetOrderItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert the ID from the URL parameter to an ObjectID
	itemID := c.Param("itemId")
	objItemID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID format"})
		return
	}

	var order models.Order
	if err := db.DB.Collection("orders").FindOne(ctx, bson.M{"items._id": objItemID}).Decode(&order); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "order item not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Iterate over the order items to find the specific item
	for _, item := range order.Items {
		if item.ItemID == objItemID {
			c.JSON(http.StatusOK, item)
			return
		}
	}

	// If the item is not found in the loop
	c.JSON(http.StatusNotFound, gin.H{"error": "order item not found"})
}
