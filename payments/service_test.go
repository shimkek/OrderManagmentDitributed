package main

import (
	"testing"

	"context"

	"github.com/shimkek/omd-common/api"
	mockRegistry "github.com/shimkek/omd-common/discovery/mock"
	"github.com/shimkek/omd-payments/gateway"
	"github.com/shimkek/omd-payments/processor/inmem"
)

func TestService(t *testing.T) {
	processor := inmem.NewInmemProcessor()
	mockRegistry := mockRegistry.NewMockRegistry()
	gateway := gateway.NewGRPCGateway(mockRegistry)
	service := NewService(processor, gateway)

	t.Run("should create payment link", func(t *testing.T) {
		link, err := service.CreatePayment(context.Background(), &api.Order{})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if link == "" {
			t.Error("expected a non-empty payment link")
		}
	})
}
