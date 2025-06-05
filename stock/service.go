package main

import (
	"context"

	"github.com/shimkek/omd-common/api"
)

type Service struct {
	store StockStore
}

func NewService(store StockStore) *Service {
	return &Service{store}
}

func (s *Service) CheckIfItemsAreInStock(ctx context.Context, p []*api.OrderItem) (bool, []*api.OrderItem, error) {
	itemIDs := make([]string, 0)
	for _, item := range p {
		itemIDs = append(itemIDs, item.ProductID)
	}

	itemsInStock, err := s.store.GetItems(ctx, itemIDs)
	if err != nil {
		return false, nil, err
	}

	// Check if all items are in stock
	for _, stockItem := range itemsInStock {
		for _, reqItem := range p {
			if stockItem.ProductID == reqItem.ProductID && stockItem.Quantity < reqItem.Quantity {
				return false, itemsInStock, nil
			}
		}
	}

	// create items with prices from stock
	items := make([]*api.OrderItem, 0)
	for _, stockItem := range itemsInStock {
		for _, reqItem := range p {
			if stockItem.ProductID == reqItem.ProductID {
				items = append(items, &api.OrderItem{
					ProductID:   stockItem.ProductID,
					ProductName: stockItem.ProductName,
					PriceID:     stockItem.PriceID,
					Quantity:    reqItem.Quantity,
				})
			}
		}
	}

	return true, items, nil
}

func (s *Service) GetItems(ctx context.Context, ids []string) ([]*api.OrderItem, error) {
	return s.store.GetItems(ctx, ids)
}
