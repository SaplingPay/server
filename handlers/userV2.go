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

const CollectionNameUserV2 = "userV2"

// CreateUserV2 creates a new user in the database
func CreateUserV2(c *gin.Context) {
	log.Println("CreateUser V2")

	var user models.UserV2

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.DB.Collection(CollectionNameUserV2).InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateUserV2 updates an existing user in the database
func UpdateUserV2(c *gin.Context) {
	log.Println("UpdateUser V2")

	userID := c.Param("userId") // Get the ID from the URL parameter

	var user models.UserV2
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	result, err := db.DB.Collection(CollectionNameUserV2).ReplaceOne(ctx, bson.M{"_id": objID}, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUserV2 deletes a user from the database
func DeleteUserV2(c *gin.Context) {
	log.Println("DeleteUser V2")

	userID := c.Param("userId")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	_, err = db.DB.Collection(CollectionNameUserV2).DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetUserV2 retrieves a single user from the database
func GetUserV2(c *gin.Context) {
	log.Println("GetUser V2")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID := c.Param("userId")

	log.Println("userID", userID)
	var user models.UserV2
	if err := db.DB.Collection(CollectionNameUserV2).FindOne(ctx, bson.M{"id": userID}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("user not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			log.Println("error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	log.Println("user", user)

	c.JSON(http.StatusOK, user)
}

// GetAllUsersV2 retrieves all users from the database
func GetAllUsersV2(c *gin.Context) {
	log.Println("GetAllUsers V2")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.UserV2
	cursor, err := db.DB.Collection(CollectionNameUserV2).Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.UserV2
		if err := cursor.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}
