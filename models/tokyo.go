package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Address   string  `bson:"address" json:"address"`
	Longitude float64 `bson:"longitude" json:"longitude"`
	Latitude  float64 `bson:"latitude" json:"latitude"`
	City      string  `bson:"city" json:"city"`
	Country   string  `bson:"country" json:"country"`
}

type Venue struct {
	ID                primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name              string               `bson:"name" json:"name"`
	Location          Location             `bson:"location" json:"location"`
	MenuID            primitive.ObjectID   `bson:"menu_id" json:"menu_id"`
	MenuIDs           []primitive.ObjectID `bson:"menu_ids" json:"menu_ids"`
	ProfilePicURL     string               `bson:"profile_pic_url" json:"profile_pic_url"`
	StripeAccountID   string               `bson:"stripe_account_id" json:"stripe_account_id"`
	OrderingSupported bool                 `bson:"ordering_supported" json:"ordering_supported"`
	DeletedAt         *primitive.DateTime  `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"` // nil if not deleted
}

type MenuV2 struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name      string              `bson:"name" json:"name"`
	VenueID   primitive.ObjectID  `bson:"venue_id" json:"venue_id"`
	Items     []MenuItemV2        `bson:"items" json:"items"`
	DeletedAt *primitive.DateTime `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"` // nil if not deleted
	// ProfileIconURL string			`bson:"profile_icon_url" json:"profile_icon_url"`
	// BannerURL string             `bson:"banner_url" json:"banner_url"`
}

type Save struct {
	Type       string             `bson:"type" json:"type"`
	VenueID    primitive.ObjectID `bson:"venue_id" json:"venue_id"`
	MenuID     primitive.ObjectID `bson:"menu_id" json:"menu_id"`
	MenuItemID primitive.ObjectID `bson:"menu_item_id" json:"menu_item_id"`
	// ADD Collections feature / field
}

type UserSavesResponse struct {
	Type          string             `bson:"type" json:"type"`
	VenueID       primitive.ObjectID `bson:"venue_id" json:"venue_id"`
	MenuID        primitive.ObjectID `bson:"menu_id" json:"menu_id"`
	MenuItemID    primitive.ObjectID `bson:"menu_item_id" json:"menu_item_id"`
	Name          string             `bson:"name" json:"name"`
	VenueName     string             `bson:"venue_name" json:"venue_name"`
	ProfilePicURL string             `bson:"profile_pic_url" json:"profile_pic_url"`
	Location      Location           `bson:"location" json:"location"`
}

type UserV2 struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID        string               `bson:"user_id" json:"user_id"`
	DisplayName   string               `bson:"display_name" json:"display_name"`
	Username      string               `bson:"username" json:"username"`
	Email         string               `bson:"email" json:"email"`
	ProfilePicURL string               `bson:"profile_pic_url" json:"profile_pic_url"`
	Location      Location             `bson:"location" json:"location"`
	Saves         []Save               `bson:"saves" json:"saves"`
	Followers     []primitive.ObjectID `bson:"followers" json:"followers"`
	Following     []primitive.ObjectID `bson:"following" json:"following"`
	DeletedAt     *primitive.DateTime  `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"` // nil if not deleted
}

type MenuItemV2 struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string              `bson:"name" json:"name"`
	Price      float64             `bson:"price" json:"price"`
	Categories []string            `bson:"categories" json:"categories"`
	DeletedAt  *primitive.DateTime `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"` // nil if not deleted
	// ADD BACK - Description, Dietary Restrictions, Ingredients, Allergens, Customizations
}

// Amounts are in cents : $1.00 = 100, $0.50 = 50, $245 = 24500
//
//	type Amount struct {
//		Cents int64 `bson:"cents" json:"cents"`
//	}
type OrderItem struct {
	MenuItemID primitive.ObjectID `bson:"menu_item_id" json:"menu_item_id"`
	Name       string             `bson:"name" json:"name"`
	Price      float64            `bson:"price" json:"price"`
	Quantity   int                `bson:"quantity" json:"quantity"`
}

type Order struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	VenueID primitive.ObjectID `bson:"venue_id" json:"venue_id"`
	// Number    int                 `bson:"number" json:"number"`
	Items     []OrderItem         `bson:"items" json:"items"`
	Total     float64             `bson:"total" json:"total"`
	Timestamp primitive.DateTime  `bson:"timestamp" json:"timestamp"`
	Status    string              `bson:"status" json:"status"`
	DeletedAt *primitive.DateTime `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"` // nil if not deleted
}

type Payment struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	OrderID   primitive.ObjectID  `bson:"order_id" json:"order_id"`
	StripeID  string              `bson:"stripe_id" json:"stripe_id"`
	Amount    float64             `bson:"amount" json:"amount"`
	Status    string              `bson:"status" json:"status"`
	Timestamp primitive.DateTime  `bson:"timestamp" json:"timestamp"`
	DeletedAt *primitive.DateTime `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"` // nil if not deleted
}

// StripeAccount temp hack figure out merchant accounts and their relation with venues
type StripeAccount struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	StripeAccountID string             `bson:"stripe_account_id" json:"stripe_account_id"`
	VenueID         primitive.ObjectID `bson:"venue_id" json:"venue_id"`
}
