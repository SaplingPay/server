package repositories

import (
	"context"
	"github.com/SaplingPay/server/db"
	"github.com/SaplingPay/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func AddMenuItem(menuId primitive.ObjectID, menuItem models.MenuItemV2) (models.MenuItemV2, error) {
	log.Println(menuItem)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println(menuId)

	// Generate a new ObjectID for the menu item
	menuItem.ID = primitive.NewObjectID()

	// Add the new menu item to the menu document
	filter := bson.M{"_id": menuId}
	update := bson.M{"$push": bson.M{"items": menuItem}}
	_, err := db.DB.Collection(db.CollectionNameMenuV2).UpdateOne(ctx, filter, update)

	return menuItem, err
}

func AddAllMenuItems(menuId primitive.ObjectID, items []models.MenuItemV2) ([]models.MenuItemV2, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set all object ids
	for idx := range items {
		items[idx].ID = primitive.NewObjectID()
	}

	filter := bson.M{"_id": menuId}
	update := bson.M{"$push": bson.M{"items": bson.M{"$each": items}}}
	_, err := db.DB.Collection(db.CollectionNameMenuV2).UpdateOne(ctx, filter, update)

	return items, err
}

func CreateMenu(menu models.MenuV2) (models.MenuV2, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.DB.Collection("menusV2").InsertOne(ctx, menu)
	if err != nil {
		return models.MenuV2{}, err
	}

	log.Println(result.InsertedID)

	venueUpdate := bson.M{"$push": bson.M{"menu_ids": menu.ID}}
	_, err = db.DB.Collection("venues").UpdateOne(ctx, bson.M{"_id": menu.VenueID}, venueUpdate)

	return menu, err
}
