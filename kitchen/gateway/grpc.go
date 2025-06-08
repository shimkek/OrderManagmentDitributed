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

func (g *gateway) UpdateOrderStatus(ctx context.Context, orderID, status string) error {
	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)
	if err != nil {
		return err
	}

	client := api.NewOrderServiceClient(conn)
	_, err = client.UpdateOrder(ctx, &api.Order{
		OrderID: orderID,
		Status:  status,
	})
	return err
}
