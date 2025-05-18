package main

import "context"

type OrderService interface {
	CreateOrder(context.Context) error
}

type OrdersStore interface {
	Create(context.Context) error
	Get(context.Context) error
	Update(context.Context) error
	Delete(context.Context) error
}
