package main

import (
	"context"

	"github.com/shimkek/omd-common/api"
)

type OrderService interface {
	CreateOrder(context.Context, *api.CreateOrderRequest) (*api.Order, error)
	ValidateOrder(context.Context, []*api.OrderItem) ([]*api.OrderItem, error)
}

type OrdersStore interface {
	Create(context.Context) error
	Get(context.Context) error
	Update(context.Context) error
	Delete(context.Context) error
}
