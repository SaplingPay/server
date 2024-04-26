package repositories

import (
	"context"
	"github.com/SaplingPay/server/db"
	"github.com/SaplingPay/server/models"
	"github.com/stripe/stripe-go/v78"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func CreatePayment(order *models.Order, session *stripe.CheckoutSession) (models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	payment := models.Payment{
		ID:       primitive.NewObjectID(),
		Amount:   float64(session.AmountTotal),
		Status:   string(session.Status),
		OrderID:  order.ID,
		StripeID: session.ID,
	}

	_, err := db.DB.Collection(db.CollectionNamePayments).InsertOne(ctx, payment)

	return payment, err
}
