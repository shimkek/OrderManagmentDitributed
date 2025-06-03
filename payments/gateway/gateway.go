package gateway

import "context"

type OrdersGateway interface {
	UpdateOrderPaymentLink(ctx context.Context, orderID, paymentLink string) error
}
