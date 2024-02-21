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

// CreatePayment creates a new payment for an order
func CreatePayment(c *gin.Context) {
	var payment models.Payment

	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.DB.Collection("payments").InsertOne(ctx, payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAllPayments retrieves all payments
func GetAllPayments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var payments []models.Payment
	cursor, err := db.DB.Collection("payments").Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var payment models.Payment
		if err := cursor.Decode(&payment); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		payments = append(payments, payment)
	}

	c.JSON(http.StatusOK, payments)
}

// GetPayment retrieves a single payment by ID
func GetPayment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paymentID := c.Param("paymentId")
	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID format"})
		return
	}

	var payment models.Payment
	if err := db.DB.Collection("payments").FindOne(ctx, bson.M{"_id": objID}).Decode(&payment); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, payment)
}

// UpdatePayment updates an existing payment
func UpdatePayment(c *gin.Context) {
	paymentID := c.Param("paymentId")

	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID format"})
		return
	}

	// Exclude the ID field from the update
	payment.ID = primitive.NilObjectID

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": payment}
	_, err = db.DB.Collection("payments").UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// DeletePayment deletes a payment
func DeletePayment(c *gin.Context) {
	paymentID := c.Param("paymentId")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID format"})
		return
	}

	filter := bson.M{"_id": objID}
	_, err = db.DB.Collection("payments").DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
}
