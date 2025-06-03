package main

import (
	"context"
	"strconv"

	common "github.com/shimkek/omd-common"
	"github.com/shimkek/omd-common/api"
)

var orders = make([]*api.Order, 0)

type store struct {
	//mongodb instance
}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(ctx context.Context, r *api.CreateOrderRequest) (*api.Order, error) {
	o := &api.Order{
		OrderID:     strconv.Itoa(len(orders) + 1),
		CustomerID:  r.CustomerID,
		Items:       r.Items,
		Status:      "pending",
		PaymentLink: "none",
	}
	orders = append(orders, o)
	return o, nil
}

func (s *store) Update(ctx context.Context, orderID string, p *api.Order) (*api.Order, error) {
	for i, order := range orders {
		if order.OrderID == orderID {
			if p.Status != "" {
				orders[i].Status = p.Status
			}
			if p.PaymentLink != "" {
				orders[i].PaymentLink = p.PaymentLink
			}
		}
	}
	return nil, nil
}
func (s *store) Get(ctx context.Context, orderID string) (*api.Order, error) {
	for _, o := range orders {
		if o.OrderID == orderID {
			return o, nil
		}
	}
	return nil, common.ErrOrderNotFound
}
func (s *store) Delete(ctx context.Context) error {
	return nil
}
