package main

import (
	"context"
	"log"

	"github.com/shimkek/omd-common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	api.UnimplementedOrderServiceServer

	service OrderService
}

func NewGrpcHandler(grpcServer *grpc.Server, service OrderService) {
	handler := &grpcHandler{service: service}
	api.RegisterOrderServiceServer(grpcServer, handler)

}

func (h *grpcHandler) CreateOrder(ctx context.Context, r *api.CreateOrderRequest) (*api.Order, error) {
	log.Println("CreateOrder gRPC handler called")
	h.service.ValidateOrder(ctx, r)
	o := &api.Order{
		OrderID: "123",
		Items:   r.Items}
	log.Printf("Order: %v", o)
	return o, nil
}
