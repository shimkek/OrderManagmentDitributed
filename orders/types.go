package main

import (
	"context"

	"github.com/shimkek/omd-common/api"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type OrderService interface {
	CreateOrder(context.Context, *api.CreateOrderRequest) (*api.Order, error)
	ValidateOrder(context.Context, []*api.OrderItem) ([]*api.OrderItem, error)
	GetOrder(context.Context, *api.GetOrderRequest) (*api.Order, error)
	UpdateOrder(context.Context, *api.Order) (*api.Order, error)
}

type OrdersStore interface {
	Create(context.Context, Order) (bson.ObjectID, error)
	Get(ctx context.Context, orderID string) (*Order, error)
	Update(context.Context, string, *api.Order) error
	Delete(context.Context) error
}

type Order struct {
	OrderID     bson.ObjectID    `bson:"_id,omitempty"`
	CustomerID  string           `bson:"customerID,omitempty"`
	Status      string           `bson:"status,omitempty"`
	PaymentLink string           `bson:"paymentLink,omitempty"`
	Items       []*api.OrderItem `bson:"items,omitempty"`
}

func (o *Order) ToProto() *api.Order {
	return &api.Order{
		OrderID:     o.OrderID.Hex(),
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
		Items:       o.Items,
	}
}
