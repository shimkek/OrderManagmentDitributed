package main

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shimkek/omd-common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	api.UnimplementedStockServiceServer

	service StockService
	channel *amqp.Channel
}

func NewGrpcHandler(grpcServer *grpc.Server, service StockService, ch *amqp.Channel) {
	handler := &grpcHandler{service: service, channel: ch}
	api.RegisterStockServiceServer(grpcServer, handler)

}

func (h *grpcHandler) CheckIfItemsAreInStock(ctx context.Context, p *api.CheckIfItemsAreInStockRequest) (*api.CheckIfItemsAreInStockResponse, error) {
	log.Printf("CheckIfItemsAreInStock called")
	inStock, items, err := h.service.CheckIfItemsAreInStock(ctx, p.Items)
	if err != nil {
		return nil, err
	}

	return &api.CheckIfItemsAreInStockResponse{
		InStock: inStock,
		Items:   items,
	}, nil
}
func (h *grpcHandler) GetItems(ctx context.Context, p *api.GetItemsRequest) (*api.GetItemsResponse, error) {
	items, err := h.service.GetItems(ctx, p.ItemIDs)
	if err != nil {
		return nil, err
	}

	return &api.GetItemsResponse{
		Items: items,
	}, nil
}
