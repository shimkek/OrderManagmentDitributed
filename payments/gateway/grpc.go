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

func (g *gateway) UpdateOrderPaymentLink(ctx context.Context, orderID, paymentLink string) error {
	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)
	if err != nil {
		return err
	}

	client := api.NewOrderServiceClient(conn)
	_, err = client.UpdateOrder(ctx, &api.Order{
		OrderID:     orderID,
		Status:      "waiting_payment",
		PaymentLink: paymentLink,
	})
	return err
}
