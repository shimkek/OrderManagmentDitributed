package gateway

import (
	"context"

	"github.com/shimkek/omd-common/api"
	"github.com/shimkek/omd-common/discovery"
)

type gateway struct {
	registry discovery.Registry
}

func NewGRPCGateway(registry discovery.Registry) *gateway {
	return &gateway{
		registry: registry,
	}
}

func (g *gateway) CreateOrder(ctx context.Context, r *api.CreateOrderRequest) (*api.Order, error) {
	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)
	if err != nil {
		return &api.Order{}, err
	}

	client := api.NewOrderServiceClient(conn)
	return client.CreateOrder(ctx, &api.CreateOrderRequest{
		CustomerID: r.CustomerID,
		Items:      r.Items,
	})
}
