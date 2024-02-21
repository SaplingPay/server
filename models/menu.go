package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MenuItem represents a single item on the digital menu.
type MenuItem struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"`
	Category    string             `bson:"category" json:"category"`
	ImageURL    string             `bson:"image_url" json:"image_url"`
	Ingredients []string           `bson:"ingredients" json:"ingredients"`
	Allergens   []string           `bson:"allergens" json:"allergens"`
}

type Menu struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Items []MenuItem         `bson:"items" json:"items"`
}

// OrderItem represents an individual item within an order.
type OrderItem struct {
	ItemID          primitive.ObjectID `bson:"item_id,omitempty" json:"item_id,omitempty"`
	Quantity        int                `bson:"quantity" json:"quantity"`
	SpecialRequests string             `bson:"special_requests" json:"special_requests"`
}

// Order represents a customer's order.
type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TableNumber int                `bson:"table_number" json:"table_number"`
	Items       []OrderItem        `bson:"items" json:"items"`
	Status      string             `bson:"status" json:"status"`
	Timestamp   time.Time          `bson:"timestamp" json:"timestamp"`
}

// Payment represents a payment made for an order.
type Payment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OrderID   primitive.ObjectID `bson:"order_id" json:"order_id"`
	Amount    float64            `bson:"amount" json:"amount"`
	Method    string             `bson:"method" json:"method"`
	Status    string             `bson:"status" json:"status"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}

// KitchenOrder represents the view of an order from the kitchen's perspective.
type KitchenOrder struct {
	OrderID     primitive.ObjectID `bson:"order_id,omitempty" json:"order_id,omitempty"`
	TableNumber int                `bson:"table_number" json:"table_number"`
	Items       []OrderItem        `bson:"items" json:"items"`
	Status      string             `bson:"status" json:"status"`
	Priority    string             `bson:"priority" json:"priority"`
}
