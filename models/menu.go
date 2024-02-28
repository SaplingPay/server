package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Items       []MenuItem         `bson:"items" json:"items"`
	UserID      string             `bson:"user_id,omitempty" json:"user_id"`
	BannerURL   string             `bson:"banner_url" json:"banner_url"`
	Location    string             `bson:"location" json:"location"`
	Description string             `bson:"description" json:"description"`
}

// MenuItem represents a single item on the digital menu.
type MenuItem struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name                string             `bson:"name" json:"name"`
	Description         string             `bson:"description" json:"description"`
	Price               float64            `bson:"price" json:"price"`
	Categories          []string           `bson:"categories" json:"categories"`
	ImageURL            string             `bson:"image_url" json:"image_url"`
	Ingredients         []string           `bson:"ingredients" json:"ingredients"`
	Allergens           []string           `bson:"allergens" json:"allergens"`
	Customizations      []string           `bson:"customizations" json:"customizations"`
	DietaryRestrictions []string           `bson:"dietary_restrictions" json:"dietary_restrictions"`
}
