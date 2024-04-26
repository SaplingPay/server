package repositories

import (
	"context"
	"github.com/SaplingPay/server/db"
	"github.com/SaplingPay/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func GetVenueByID(venueID primitive.ObjectID) (models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var order models.Order
	err := db.DB.Collection(db.CollectionNameVenue).FindOne(ctx, bson.M{"_id": venueID}).Decode(&order)

	return order, err
}
