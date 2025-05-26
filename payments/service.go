package main

import (
	"context"
	"log"

	"github.com/shimkek/omd-common/api"
	"github.com/shimkek/omd-payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
}

func NewService(processor processor.PaymentProcessor) *service {
	return &service{processor: processor}
}
func (s *service) CreatePayment(ctx context.Context, order *api.Order) (string, error) {
	link, err := s.processor.CreatePaymentLink(order)
	if err != nil {
		log.Printf("failed to create payment link: %v", err)
		return "", err
	}
	return link, nil
}
