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

const CollectionNameVenue = "venues"

// CreateVenue creates a new venue in the database
func CreateVenue(c *gin.Context) {
	log.Println("CreateVenue")

	var venue models.Venue

	if err := c.ShouldBindJSON(&venue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	venue.MenuIDs = []primitive.ObjectID{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.DB.Collection(CollectionNameVenue).InsertOne(ctx, venue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Venue created:", result.InsertedID)

	c.JSON(http.StatusOK, result)
}

// UpdateVenue updates an existing venue in the database
func UpdateVenue(c *gin.Context) {
	log.Println("UpdateVenue")

	venueID := c.Param("venueId") // Get the ID from the URL parameter

	var venue models.Venue
	if err := c.ShouldBindJSON(&venue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert the string ID to MongoDB's ObjectID
	objID, err := primitive.ObjectIDFromHex(venueID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	// Prepare the update document similarly to before
	update := bson.M{}
	venueType := reflect.TypeOf(venue)
	venueValue := reflect.ValueOf(venue)
	for i := 0; i < venueType.NumField(); i++ {
		field := venueType.Field(i)
		fieldValue := venueValue.Field(i).Interface()
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

	_, err = db.DB.Collection(CollectionNameVenue).UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the updated venue from the database
	var updatedVenue models.Venue
	err = db.DB.Collection(CollectionNameVenue).FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedVenue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve updated venue"})
		return
	}

	c.JSON(http.StatusOK, updatedVenue)
}

// DeleteVenue deletes a venue from the database
func DeleteVenue(c *gin.Context) {
	log.Println("DeleteVenue")

	// Fetching the venue ID from the URL parameter
	venueID := c.Param("venueId")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Assuming `venueID` needs to be converted to an ObjectID if you're using MongoDB's default ObjectID
	// If your ID is a string in the database, you can directly use venueID in the filter
	objID, err := primitive.ObjectIDFromHex(venueID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	_, err = db.DB.Collection(CollectionNameVenue).DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Venue deleted successfully"})
}

// GetVenue retrieves a single venue from the database
func GetVenue(c *gin.Context) {
	log.Println("GetVenue")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert the ID from the URL parameter to an ObjectID
	venueID := c.Param("venueId")
	objID, err := primitive.ObjectIDFromHex(venueID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var venue models.Venue
	// Use the ObjectID to find the venue
	if err := db.DB.Collection(CollectionNameVenue).FindOne(ctx, bson.M{"_id": objID}).Decode(&venue); err != nil {
		// Adjust the error handling to distinguish not found errors from other errors
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "venue not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, venue)
}

// GetAllVenues retrieves all venues from the database
func GetAllVenues(c *gin.Context) {
	log.Println("GetAllVenues")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var venues []models.Venue
	cursor, err := db.DB.Collection(CollectionNameVenue).Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var venue models.Venue
		if err := cursor.Decode(&venue); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		venues = append(venues, venue)
	}

	c.JSON(http.StatusOK, venues)
}
