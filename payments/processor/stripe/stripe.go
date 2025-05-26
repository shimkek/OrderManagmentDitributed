package StripeProcessor

import (
	"fmt"
	"log"

	common "github.com/shimkek/omd-common"
	"github.com/shimkek/omd-common/api"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

var (
	gatewayHTTPAddr   = common.EnvGetString("GATEWAY_HTTP_ADDR", "http://localhost:8080")
	gatewaySuccessURL = fmt.Sprintf("%s/success.html", gatewayHTTPAddr)
)

type StripeProcessor struct{}

func NewProcessor() *StripeProcessor {
	return &StripeProcessor{}
}
func (s *StripeProcessor) CreatePaymentLink(order *api.Order) (string, error) {
	log.Printf("Creating payment link for order: %s", order.OrderID)

	items := make([]*stripe.CheckoutSessionLineItemParams, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String("price_1RT6Lg07x6KdXnTOdZW6J5rc"),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}
	checkoutParams := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems:  items,
		SuccessURL: stripe.String(gatewaySuccessURL),
	}
	result, err := session.New(checkoutParams)
	if err != nil {
		log.Printf("Failed to create checkout session: %v", err)
		return "", err
	}
	return result.URL, nil

}
