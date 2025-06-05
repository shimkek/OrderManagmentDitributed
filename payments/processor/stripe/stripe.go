package StripeProcessor

import (
	"fmt"
	"log"

	common "github.com/shimkek/omd-common"
	"github.com/shimkek/omd-common/api"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

var (
	gatewayHTTPAddr = common.EnvGetString("GATEWAY_HTTP_ADDR", "http://localhost:8080")
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
			Price:    stripe.String(item.PriceID),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}
	gatewaySuccessURL := fmt.Sprintf("%s/success.html?customerID=%s&orderID=%s", gatewayHTTPAddr, order.CustomerID, order.OrderID)
	checkoutParams := &stripe.CheckoutSessionParams{
		Metadata: map[string]string{
			"OrderID":    order.OrderID,
			"CustomerID": order.CustomerID,
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems:  items,
		SuccessURL: stripe.String(gatewaySuccessURL),
		CancelURL:  stripe.String(fmt.Sprintf("%s/cancel.html", gatewayHTTPAddr)),
	}
	result, err := session.New(checkoutParams)
	if err != nil {
		log.Printf("Failed to create checkout session: %v", err)
		return "", err
	}
	return result.URL, nil

}
