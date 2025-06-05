package gateway

import (
	"context"

	"github.com/shimkek/omd-common/api"
)

type StockGateway interface {
	CheckIfItemsAreInStock(ctx context.Context, items []*api.OrderItem) (bool, []*api.OrderItem, error)
	GetItems(ctx context.Context, p *api.GetItemsRequest) (*api.GetItemsResponse, error)
}
