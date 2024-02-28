package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Payment represents a payment made for an order.
type Payment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OrderID   primitive.ObjectID `bson:"order_id" json:"order_id"`
	Amount    float64            `bson:"amount" json:"amount"`
	Method    string             `bson:"method" json:"method"`
	Status    string             `bson:"status" json:"status"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}
