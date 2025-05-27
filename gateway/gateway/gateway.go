package gateway

import (
	"context"

	"github.com/shimkek/omd-common/api"
)

type OrdersGateway interface {
	CreateOrder(context.Context, *api.CreateOrderRequest) (*api.Order, error)
	GetOrder(context.Context, string) (*api.Order, error)
}
