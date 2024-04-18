package handlers

import (
	"context"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/SaplingPay/server/db"
	"github.com/SaplingPay/server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const CollectionNameUserV2 = "usersV2"

// CreateUserV2 creates a new user in the database
func CreateUserV2(c *gin.Context) {
	log.Println("CreateUser V2")

	var user models.UserV2

	// init
	user.Followers = []primitive.ObjectID{}
	user.Following = []primitive.ObjectID{}
	user.Saves = []models.Save{}

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

	// Prepare the update document similarly to before
	update := bson.M{}
	userType := reflect.TypeOf(user)
	userValue := reflect.ValueOf(user)
	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)
		fieldValue := userValue.Field(i).Interface()
		fieldType := field.Type.Kind()

		if fieldType == reflect.Bool || !reflect.DeepEqual(fieldValue, reflect.Zero(field.Type).Interface()) {
			bsonTag := field.Tag.Get("bson")
			// Skip if bson tag is not set or is "-"
			if bsonTag == "" || bsonTag == "-" {
				continue
			}

			update[field.Tag.Get("bson")] = fieldValue
		}
	}

	result, err := db.DB.Collection(CollectionNameUserV2).UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Retrieve the updated user from the database
	var updatedUser models.UserV2
	err = db.DB.Collection(CollectionNameUserV2).FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve updated user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
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
	if err := db.DB.Collection(CollectionNameUserV2).FindOne(ctx, bson.M{"user_id": userID}).Decode(&user); err != nil {
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

func FollowUser(c *gin.Context) {
	log.Println("Follow")

	userID := c.Param("userId")
	followingID := c.Param("followingId")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	followingObjID, err := primitive.ObjectIDFromHex(followingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	_, err = db.DB.Collection(CollectionNameUserV2).UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$addToSet": bson.M{"following": followingObjID.Hex()}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = db.DB.Collection(CollectionNameUserV2).UpdateOne(ctx, bson.M{"_id": followingObjID}, bson.M{"$addToSet": bson.M{"followers": objID.Hex()}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updatedUser models.UserV2
	err = db.DB.Collection(CollectionNameUserV2).FindOne(ctx, bson.M{"_id": followingObjID}).Decode(&updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve updated user"})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

func UnFollowUser(c *gin.Context) {
	log.Println("UnFollow")

	userID := c.Param("userId")
	followingID := c.Param("followingId")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	followingObjID, err := primitive.ObjectIDFromHex(followingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	_, err = db.DB.Collection(CollectionNameUserV2).UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$pull": bson.M{"following": followingObjID.Hex()}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = db.DB.Collection(CollectionNameUserV2).UpdateOne(ctx, bson.M{"_id": followingObjID}, bson.M{"$pull": bson.M{"followers": objID.Hex()}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updatedUser models.UserV2
	err = db.DB.Collection(CollectionNameUserV2).FindOne(ctx, bson.M{"_id": followingObjID}).Decode(&updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve updated user"})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

// GetUserSaves retrieves all saves for a user
func GetUserSaves(c *gin.Context) {
	log.Println("GetUserSaves")

	userID := c.Param("userId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var user models.UserV2
	err = db.DB.Collection(CollectionNameUserV2).FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var userSaves []models.UserSavesResponse
	for _, save := range user.Saves {

		var userSave models.UserSavesResponse
		if save.Type == "venue" {
			var venue models.Venue
			err = db.DB.Collection("venues").FindOne(ctx, bson.M{"_id": save.VenueID}).Decode(&venue)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve venue"})
				return
			}

			userSave.Type = save.Type
			userSave.VenueID = save.VenueID
			userSave.MenuID = save.MenuID
			userSave.MenuItemID = save.MenuItemID
			userSave.Name = venue.Name
			userSave.VenueName = venue.Name
			userSave.ProfilePicURL = venue.ProfilePicURL
			userSave.Location = venue.Location
		} else if save.Type == "menu_item" {
			var menu models.MenuV2
			err = db.DB.Collection("menusV2").FindOne(ctx, bson.M{"_id": save.MenuID}).Decode(&menu)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve menu"})
				return
			}

			var venue models.Venue
			err = db.DB.Collection("venues").FindOne(ctx, bson.M{"_id": menu.VenueID}).Decode(&venue)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve venue"})
				return
			}

			var menuItem models.MenuItemV2
			for _, item := range menu.Items {
				if item.ID == save.MenuItemID {
					menuItem = item
					break
				}
			}

			userSave.Type = save.Type
			userSave.VenueID = save.VenueID
			userSave.MenuID = save.MenuID
			userSave.MenuItemID = save.MenuItemID
			userSave.Name = menuItem.Name
			userSave.VenueName = venue.Name
			userSave.ProfilePicURL = venue.ProfilePicURL
			userSave.Location = venue.Location
		}

		userSaves = append(userSaves, userSave)
	}

	c.JSON(http.StatusOK, userSaves)
}
