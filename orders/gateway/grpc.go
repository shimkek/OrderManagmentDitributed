package gateway

import (
	"context"

	"github.com/shimkek/omd-common/api"
	"github.com/shimkek/omd-common/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) CheckIfItemsAreInStock(ctx context.Context, items []*api.OrderItem) (bool, []*api.OrderItem, error) {
	conn, err := discovery.ServiceConnection(ctx, "stock", g.registry)
	if err != nil {
		return false, nil, err
	}
	defer conn.Close()

	client := api.NewStockServiceClient(conn)
	res, err := client.CheckIfItemsAreInStock(ctx, &api.CheckIfItemsAreInStockRequest{
		Items: items,
	})
	if err != nil {
		return false, nil, err
	}
	return res.InStock, res.Items, nil
}
func (h *Gateway) GetItems(ctx context.Context, p *api.GetItemsRequest) (*api.GetItemsResponse, error) {
	return nil, nil
}
