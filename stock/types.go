package main

import (
	"context"

	"github.com/shimkek/omd-common/api"
)

type StockService interface {
	CheckIfItemsAreInStock(context.Context, []*api.OrderItem) (bool, []*api.OrderItem, error)
	GetItems(ctx context.Context, ids []string) ([]*api.OrderItem, error)
}

type StockStore interface {
	GetItem(ctx context.Context, id string) (*api.OrderItem, error)
	GetItems(ctx context.Context, ids []string) ([]*api.OrderItem, error)
}
