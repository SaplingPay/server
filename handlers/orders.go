package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/SaplingPay/server/repositories"

	"github.com/SaplingPay/server/db"
	"github.com/SaplingPay/server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.ID = primitive.NewObjectID() // Generate a new ID for the order
	order.Timestamp = primitive.NewDateTimeFromTime(time.Now())
	order.Status = "sent" // Set the default status

	// Assuming there's logic to calculate the total from order.Items
	order.Total = calculateTotal(order.Items)

	_, err := db.DB.Collection(db.CollectionNameOrders).InsertOne(context.Background(), order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	order, err := repositories.GetOrderByID(objID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func UpdateOrder(c *gin.Context) {
	orderID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if items, exists := updates["items"]; exists {
		// Calculate new total if items are updated
		updates["total"] = calculateTotal(items.([]models.OrderItem))
	}

	_, err = db.DB.Collection(db.CollectionNameOrders).UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": updates})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order updated"})
}

// func HardDeleteOrder(c *gin.Context) {
// 	orderID := c.Param("id")
// 	objID, err := primitive.ObjectIDFromHex(orderID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
// 		return
// 	}

// 	_, err = db.DB.Collection(CollectionNameOrders).DeleteOne(context.Background(), bson.M{"_id": objID})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
// }

func SoftDeleteOrder(c *gin.Context) {
	orderID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	update := bson.M{
		"$set": bson.M{"deleted_at": primitive.NewDateTimeFromTime(time.Now())},
	}

	_, err = db.DB.Collection("orders").UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order soft deleted"})
}

func GetAllOrders(c *gin.Context) {
	var orders []models.Order
	cursor, err := db.DB.Collection(db.CollectionNameOrders).Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, orders)
}

func calculateTotal(items []models.OrderItem) float64 {
	var total float64
	// Calculation logic based on items
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}
