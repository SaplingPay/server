package repositories

import (
	"context"
	"github.com/SaplingPay/server/db"
	"github.com/SaplingPay/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func GetStripeAccounts() ([]models.StripeAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var accounts []models.StripeAccount

	cursor, err := db.DB.Collection(db.CollectionNameStripeAccounts).Find(ctx, bson.M{})
	if err != nil {
		return accounts, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var account models.StripeAccount
		if err := cursor.Decode(&account); err != nil {
			return accounts, err
		}
		accounts = append(accounts, account)
	}

	return accounts, err
}

func AddStripeAccount(stripeAccountNumber string) (models.StripeAccount, error) {
	objId := primitive.NewObjectID()
	account := models.StripeAccount{
		ID:              objId,
		StripeAccountID: stripeAccountNumber,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.DB.Collection(db.CollectionNameStripeAccounts).InsertOne(ctx, &account)

	return account, err
}

func GetStripeAccountByVenueId(venueId primitive.ObjectID) (models.StripeAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var account models.StripeAccount
	err := db.DB.Collection(db.CollectionNameOrders).FindOne(ctx, bson.M{"venueId": venueId}).Decode(&account)

	return account, err
}

func LinkVenue(venueId primitive.ObjectID, accountNumber string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{"venueId": venueId},
	}

	_, err := db.DB.Collection(db.CollectionNameStripeAccounts).UpdateOne(ctx, &bson.M{"stripe_account_id": accountNumber}, update)

	return err
}
