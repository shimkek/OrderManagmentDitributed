package main

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shimkek/omd-common/api"
	"github.com/shimkek/omd-common/broker"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	api.UnimplementedOrderServiceServer

	service OrderService
	channel *amqp.Channel
}

func NewGrpcHandler(grpcServer *grpc.Server, service OrderService, ch *amqp.Channel) {
	handler := &grpcHandler{service: service, channel: ch}
	api.RegisterOrderServiceServer(grpcServer, handler)

}

func (h *grpcHandler) GetOrder(ctx context.Context, r *api.GetOrderRequest) (*api.Order, error) {
	log.Println("GetOrder gRPC handler called. OrderID: ", r.OrderID)
	o, err := h.service.GetOrder(ctx, r)
	if err != nil {
		log.Printf("Failed to get order: %v", err)
		return nil, err
	}
	log.Printf("Order:\n %v", o)

	return o, nil
}

func (h *grpcHandler) CreateOrder(ctx context.Context, r *api.CreateOrderRequest) (*api.Order, error) {
	log.Println("CreateOrder gRPC handler called")
	o, err := h.service.CreateOrder(ctx, r)
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		return nil, err
	}
	log.Printf("Order:\n %v", o)

	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		log.Printf("Failed to marshal order: %v", err)
		return nil, err
	}

	q, err := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
		return nil, err
	}

	h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshalledOrder,
		DeliveryMode: amqp.Persistent,
	})

	return o, nil
}
