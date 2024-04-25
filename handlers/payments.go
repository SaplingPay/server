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
)

func CreatePayment(c *gin.Context) {
	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment.ID = primitive.NewObjectID()                          // Generate a new ID for the payment
	payment.Timestamp = primitive.NewDateTimeFromTime(time.Now()) // Set the current timestamp

	_, err := db.DB.Collection("payments").InsertOne(context.Background(), payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func GetPayment(c *gin.Context) {
	paymentID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var payment models.Payment
	filter := bson.M{"_id": objID, "deleted_at": bson.M{"$exists": false}}
	if err := db.DB.Collection("payments").FindOne(context.Background(), filter).Decode(&payment); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func UpdatePayment(c *gin.Context) {
	paymentID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var updates bson.M
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.DB.Collection("payments").UpdateOne(context.Background(), bson.M{"_id": objID, "deleted_at": bson.M{"$exists": false}}, bson.M{"$set": updates})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment updated"})
}

func SoftDeletePayment(c *gin.Context) {
	paymentID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	update := bson.M{
		"$set": bson.M{"deleted_at": primitive.NewDateTimeFromTime(time.Now())},
	}

	_, err = db.DB.Collection("payments").UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment soft deleted"})
}

func GetAllPayments(c *gin.Context) {
	var payments []models.Payment
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}
	cursor, err := db.DB.Collection("payments").Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var payment models.Payment
		if err := cursor.Decode(&payment); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		payments = append(payments, payment)
	}

	c.JSON(http.StatusOK, payments)
}
