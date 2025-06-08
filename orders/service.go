package main

import (
	"context"

	common "github.com/shimkek/omd-common"
	"github.com/shimkek/omd-common/api"
	"github.com/shimkek/omd-orders/gateway"
)

type service struct {
	store   OrdersStore
	gateway gateway.StockGateway
}

func NewService(store OrdersStore, gateway gateway.StockGateway) *service {
	return &service{
		store:   store,
		gateway: gateway,
	}
}

func (s *service) CreateOrder(ctx context.Context, p *api.CreateOrderRequest) (*api.Order, error) {
	validatedItems, err := s.ValidateOrder(ctx, p.Items)
	if err != nil {
		return nil, err
	}
	id, err := s.store.Create(ctx, Order{
		CustomerID:  p.CustomerID,
		Items:       validatedItems,
		Status:      "pending",
		PaymentLink: "none",
	})
	if err != nil {
		return nil, err
	}
	return &api.Order{
		OrderID:     id.Hex(),
		CustomerID:  p.CustomerID,
		Items:       validatedItems,
		Status:      "pending",
		PaymentLink: "none",
	}, nil
}

func (s *service) GetOrder(ctx context.Context, r *api.GetOrderRequest) (*api.Order, error) {
	order, err := s.store.Get(ctx, r.OrderID)
	if err != nil {
		return nil, err
	}
	return order.ToProto(), nil

}

func (s *service) UpdateOrder(ctx context.Context, o *api.Order) (*api.Order, error) {
	s.store.Update(ctx, o.OrderID, o)
	return o, nil
}

func (s *service) ValidateOrder(ctx context.Context, items []*api.OrderItem) ([]*api.OrderItem, error) {
	if len(items) == 0 {
		return nil, common.ErrNoItems
	}
	mergedItems := mergeOrderItems(items)

	inStock, items, err := s.gateway.CheckIfItemsAreInStock(ctx, mergedItems)
	if err != nil {
		return nil, err
	}
	if !inStock {
		return nil, common.ErrNotInStock
	}

	return items, nil
}

func mergeOrderItems(items []*api.OrderItem) []*api.OrderItem {
	merged := make([]*api.OrderItem, 0)

	for _, item := range items {
		found := false
		for _, finalItem := range merged {
			if item.ProductID == finalItem.ProductID {
				finalItem.Quantity += item.Quantity
				found = true
				break
			}
		}
		if !found {
			merged = append(merged, item)
		}
	}
	return merged
}
