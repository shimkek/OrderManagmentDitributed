package main

import (
	"context"
	"log"

	"github.com/shimkek/omd-common/api"
	"github.com/shimkek/omd-payments/gateway"
	"github.com/shimkek/omd-payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
	gateway   gateway.OrdersGateway
}

func NewService(processor processor.PaymentProcessor, gateway gateway.OrdersGateway) *service {
	return &service{processor: processor, gateway: gateway}
}
func (s *service) CreatePayment(ctx context.Context, order *api.Order) (string, error) {
	link, err := s.processor.CreatePaymentLink(order)
	if err != nil {
		log.Printf("failed to create payment link: %v", err)
		return "", err
	}

	err = s.gateway.UpdateOrderPaymentLink(ctx, order.OrderID, link)
	if err != nil {
		return "", err
	}
	return link, nil
}
