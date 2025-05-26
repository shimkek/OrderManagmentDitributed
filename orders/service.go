package main

import (
	"context"

	common "github.com/shimkek/omd-common"
	"github.com/shimkek/omd-common/api"
)

type service struct {
	store OrdersStore
}

func NewService(store OrdersStore) *service {
	return &service{
		store: store,
	}
}

func (s *service) CreateOrder(ctx context.Context, p *api.CreateOrderRequest) (*api.Order, error) {
	validatedItems, err := s.ValidateOrder(ctx, p.Items)
	if err != nil {
		return nil, err
	}
	o := &api.Order{
		OrderID:    "123",
		CustomerID: p.CustomerID,
		Items:      validatedItems,
		Status:     "pending",
	}
	return o, nil
}

func (s *service) ValidateOrder(ctx context.Context, items []*api.OrderItem) ([]*api.OrderItem, error) {
	if len(items) == 0 {
		return nil, common.ErrNoItems
	}
	validatedItems := mergeOrderItems(items)
	return validatedItems, nil
	// validate with stock service

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
