package main

import (
	"context"
	"log"

	pb "github.com/shimkek/omd-common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
}

func NewGrpcHandler(grpcServer *grpc.Server) {
	handler := &grpcHandler{}
	pb.RegisterOrderServiceServer(grpcServer, handler)

}

func (h *grpcHandler) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Println("CreateOrder gRPC handler called")
	log.Printf("Order items: %v", r.Items)
	o := &pb.Order{
		OrderId: "123"}
	return o, nil
}
