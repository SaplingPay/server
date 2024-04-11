package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Address   string  `bson:"address" json:"address"`
	Longitude float64 `bson:"longitude" json:"longitude"`
	Latitude  float64 `bson:"latitude" json:"latitude"`
}

type Venue struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Location Location           `bson:"location" json:"location"`
	Menu     MenuV2             `bson:"menus" json:"menus"`
}

type MenuV2 struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Items     []MenuItemV2       `bson:"items" json:"items"`
	BannerURL string             `bson:"banner_url" json:"banner_url"`
	Location  string             `bson:"location" json:"location"`
}

type UserV2 struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId      string             `bson:"user_id" json:"user_id"`
	DisplayName string             `bson:"display_name" json:"display_name"`
	Username    string             `bson:"username" json:"username"`
	Email       string             `bson:"email" json:"email"`
}

type MenuItemV2 struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Price      float64            `bson:"price" json:"price"`
	Categories []string           `bson:"categories" json:"categories"`
	// ADD BACK - Description, Dietary Restrictions, Ingredients, Allergens, Customizations
}
