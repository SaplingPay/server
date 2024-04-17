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

const CollectionNameMenuV2 = "menusV2"

func CreateMenuV2(c *gin.Context) {
	log.Println("CreateMenu V2")

	venueID := c.Param("venueId") // Assuming the venue ID is passed as a URL parameter

	var menu models.MenuV2
	if err := c.ShouldBindJSON(&menu); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(menu)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert venueID from string to primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(venueID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu ID format"})
		return
	}

	log.Println(objID)

	// Generate a new ObjectID for the menu item
	menu.ID = primitive.NewObjectID()

	if menu.Items == nil {
		menu.Items = []models.MenuItemV2{}
	}

	result, err := db.DB.Collection("menusV2").InsertOne(ctx, menu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println(result.InsertedID)

	// Store the menu ID in the venue's menu ID field
	venueUpdate := bson.M{"$push": bson.M{"menu_ids": menu.ID}}
	_, err = db.DB.Collection("venues").UpdateOne(ctx, bson.M{"_id": objID}, venueUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// UpdateMenu updates an existing menu in the database
func UpdateMenuV2(c *gin.Context) {
	log.Println("UpdateMenu V2")

	menuID := c.Param("menuId") // Get the ID from the URL parameter

	var menu models.MenuV2
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

	// Prepare the update document similarly to before
	update := bson.M{}
	menuType := reflect.TypeOf(menu)
	menuValue := reflect.ValueOf(menu)
	for i := 0; i < menuType.NumField(); i++ {
		field := menuType.Field(i)
		fieldValue := menuValue.Field(i).Interface()
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

	_, err = db.DB.Collection(CollectionNameMenuV2).UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the updated menu from the database
	var updatedMenu models.MenuV2
	err = db.DB.Collection(CollectionNameMenuV2).FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedMenu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve updated menu"})
		return
	}

	c.JSON(http.StatusOK, updatedMenu)
}

// DeleteMenu deletes a menu from the database
func DeleteMenuV2(c *gin.Context) {
	log.Println("DeleteMenu V2")

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

	_, err = db.DB.Collection(CollectionNameMenuV2).DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu deleted successfully"})
}

// GetMenu retrieves a single menu from the database
func GetMenuV2(c *gin.Context) {
	log.Println("GetMenu V2")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert the ID from the URL parameter to an ObjectID
	menuID := c.Param("menuId")
	objID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var menu models.MenuV2
	// Use the ObjectID to find the menu
	if err := db.DB.Collection(CollectionNameMenuV2).FindOne(ctx, bson.M{"_id": objID}).Decode(&menu); err != nil {
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
func GetAllMenusV2(c *gin.Context) {
	log.Println("GetAllMenus V2")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var menus []models.MenuV2
	cursor, err := db.DB.Collection(CollectionNameMenuV2).Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var menu models.MenuV2
		if err := cursor.Decode(&menu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		menus = append(menus, menu)
	}

	c.JSON(http.StatusOK, menus)
}
