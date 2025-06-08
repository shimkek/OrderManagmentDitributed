package gateway

import "context"

type OrdersGateway interface {
	UpdateOrderStatus(context.Context, string, string) error
}
