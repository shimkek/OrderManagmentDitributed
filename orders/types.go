package main

import (
	"context"

	"github.com/shimkek/omd-common/api"
)

type OrderService interface {
	CreateOrder(context.Context, *api.CreateOrderRequest) (*api.Order, error)
	ValidateOrder(context.Context, []*api.OrderItem) ([]*api.OrderItem, error)
	GetOrder(context.Context, *api.GetOrderRequest) (*api.Order, error)
	UpdateOrder(context.Context, *api.Order) (*api.Order, error)
}

type OrdersStore interface {
	Create(context.Context, *api.CreateOrderRequest) (*api.Order, error)
	Get(ctx context.Context, orderID string) (*api.Order, error)
	Update(context.Context, string, *api.Order) (*api.Order, error)
	Delete(context.Context) error
}
