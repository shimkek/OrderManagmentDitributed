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

func (s *service) CreateOrder(ctx context.Context) error {
	return nil
}

func (s *service) ValidateOrder(ctx context.Context, p *api.CreateOrderRequest) error {
	if len(p.Items) == 0 {
		return common.ErrNoItems
	}
	p.Items = mergeOrderItems(p.Items)

	// validate with stock service

	return nil
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
