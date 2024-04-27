package payments

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SaplingPay/server/models"
	"github.com/SaplingPay/server/repositories"
	"github.com/SaplingPay/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/account"
	"github.com/stripe/stripe-go/v78/accountsession"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddStripRoutes(r *gin.Engine) {
	stripeRoutes := r.Group("/payments")
	{
		stripeRoutes.GET("/accounts", GetAccounts)

		stripeRoutes.POST("/linkAccount", LinkAccount)
		stripeRoutes.POST("/account", CreateAccount)
		stripeRoutes.POST("/accountSession", CreateAccountSession)
		stripeRoutes.POST("/checkout/:orderId", CreateCheckoutSession)
	}
}
func GetAccounts(c *gin.Context) {
	accounts, err := repositories.GetStripeAccounts()

	if err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": "unable to fetch stripe accounts"})
		return
	}
	c.JSON(http.StatusOK, &gin.H{"accounts": accounts})
}

func CreateAccountSession(c *gin.Context) {
	log.Println("[stripeHandler]", "CreateAccountSession")

	type RequestBody struct {
		Account string `json:"account"`
	}

	var requestBody RequestBody
	err := json.NewDecoder(c.Request.Body).Decode(&requestBody)

	if err != nil {
		c.JSON(http.StatusBadRequest, &gin.H{"error": err.Error()})
		return
	}

	params := &stripe.AccountSessionParams{
		Account: stripe.String(requestBody.Account),
		Components: &stripe.AccountSessionComponentsParams{
			AccountOnboarding: &stripe.AccountSessionComponentsAccountOnboardingParams{
				Enabled: stripe.Bool(true),
			},
			Payments: &stripe.AccountSessionComponentsPaymentsParams{
				Enabled: stripe.Bool(true),
				Features: &stripe.AccountSessionComponentsPaymentsFeaturesParams{
					RefundManagement:                      stripe.Bool(true),
					DisputeManagement:                     stripe.Bool(true),
					CapturePayments:                       stripe.Bool(true),
					DestinationOnBehalfOfChargeManagement: stripe.Bool(false),
				},
			},
		},
	}

	accountSession, err := accountsession.New(params)

	if err != nil {
		log.Printf("An error occurred when calling the Stripe API to create an account session: %v", err)
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, &gin.H{"client_secret": accountSession.ClientSecret})
}

func CreateAccount(c *gin.Context) {
	account, err := account.New(&stripe.AccountParams{
		Controller: &stripe.AccountControllerParams{
			StripeDashboard: &stripe.AccountControllerStripeDashboardParams{
				Type: stripe.String("none"),
			},
			Fees: &stripe.AccountControllerFeesParams{
				Payer: stripe.String("application"),
			},
		},
		Capabilities: &stripe.AccountCapabilitiesParams{
			CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
				Requested: stripe.Bool(true),
			},
			Transfers: &stripe.AccountCapabilitiesTransfersParams{
				Requested: stripe.Bool(true),
			},
		},
		Country: stripe.String("NL"),
	})

	if err != nil {
		log.Printf("An error occurred when calling the Stripe API to create an account: %v", err)
		handleError(c, err)
		return
	}

	_, err = repositories.AddStripeAccount(account.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": "unable to save account"})
		return
	}

	c.JSON(http.StatusOK, &gin.H{"account": account.ID})
}

func LinkAccount(c *gin.Context) {
	var requestBody models.StripeAccount
	err := json.NewDecoder(c.Request.Body).Decode(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, &gin.H{"error": err.Error()})
		return
	}

	err = repositories.LinkVenue(requestBody.VenueID, requestBody.StripeAccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error linking account"})
		return
	}

	c.JSON(http.StatusOK, requestBody)
}

func CreateCheckoutSession(c *gin.Context) {
	orderIdParam := c.Param("orderId")

	orderId, err := primitive.ObjectIDFromHex(orderIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	order, err := repositories.GetOrderByID(orderId)
	log.Println(order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID"})
		return
	}

	// _, err = repositories.GetStripeAccountByVenueId(order.VenueID)

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "No Stripe Account Linked"})
	// 	return
	// }

	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, item := range order.Items {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String(string(stripe.CurrencyEUR)),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(item.Name),
				},
				UnitAmount: stripe.Int64(int64(item.Price * 100)),
			},
			// TODO Handle quantity better
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	successURL := os.Getenv("STRIPE_SUCCESS_URL_ORIGIN")
	if successURL == "" {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": "missing success URL origin"})
		return
	}

	params := &stripe.CheckoutSessionParams{
		LineItems: lineItems,
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			ApplicationFeeAmount: stripe.Int64(0),
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(fmt.Sprintf("%s/order-received?order_id=%s", successURL, orderId.Hex())),
	}
	// todo change this to prod id acct_1P9t8BKrxYc2JpQl
	params.SetStripeAccount("acct_1P9xWfJcPRqoOyDQ")
	result, err := session.New(params)
	if err != nil {
		handleError(c, err)
		return
	}

	_, err = repositories.CreatePayment(&order, result)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorJson("error creating payment"))
	}

	// if redirect doesn't work with frontend switch to json
	c.JSON(http.StatusOK, &gin.H{"url": result.URL})
	// c.Redirect(http.StatusFound, result.URL)
}

func handleError(c *gin.Context, err error) {
	if stripeErr, ok := err.(*stripe.Error); ok {
		c.JSON(http.StatusInternalServerError, &gin.H{
			"error": stripeErr.Msg,
		})
	} else {
		c.JSON(http.StatusInternalServerError, &gin.H{
			"error": err.Error(),
		})
	}
}
