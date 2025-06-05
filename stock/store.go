package main

import (
	"context"
	"fmt"

	"github.com/shimkek/omd-common/api"
)

type Store struct {
	stock map[string]*api.OrderItem
}

func NewStore() *Store {
	return &Store{
		stock: map[string]*api.OrderItem{
			"1": {
				ProductID:   "1",
				ProductName: "Cheese",
				Quantity:    12,
				PriceID:     "price_1RWg6C07x6KdXnTOR532hPIw",
			},
			"2": {
				ProductID:   "2",
				ProductName: "Chocolate",
				Quantity:    7,
				PriceID:     "price_1RT6Lg07x6KdXnTOdZW6J5rc",
			},
		},
	}
}

func (s *Store) GetItem(ctx context.Context, id string) (*api.OrderItem, error) {
	for _, item := range s.stock {
		if item.ProductID == id {
			return item, nil
		}
	}
	return nil, fmt.Errorf("item not found")
}
func (s *Store) GetItems(ctx context.Context, ids []string) ([]*api.OrderItem, error) {
	var res []*api.OrderItem
	for _, id := range ids {
		if i, ok := s.stock[id]; ok {
			res = append(res, i)
		}
	}

	return res, nil
}
