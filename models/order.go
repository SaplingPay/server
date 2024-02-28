package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

// KitchenOrder represents the view of an order from the kitchen's perspective.
type KitchenOrder struct {
	OrderID     primitive.ObjectID `bson:"order_id,omitempty" json:"order_id,omitempty"`
	TableNumber int                `bson:"table_number" json:"table_number"`
	Items       []OrderItem        `bson:"items" json:"items"`
	Status      string             `bson:"status" json:"status"`
	Priority    string             `bson:"priority" json:"priority"`
}
