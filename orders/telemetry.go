package main

import (
	"context"
	"fmt"

	"github.com/shimkek/omd-common/api"
	"go.opentelemetry.io/otel/trace"
)

type TelemetryMiddleware struct {
	next OrderService
}

func NewTelemetryMiddleware(next OrderService) OrderService {
	return &TelemetryMiddleware{next: next}
}

func (s *TelemetryMiddleware) CreateOrder(ctx context.Context, p *api.CreateOrderRequest) (*api.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("CreateORder: %v", p))
	return s.next.CreateOrder(ctx, p)
}

func (s *TelemetryMiddleware) GetOrder(ctx context.Context, r *api.GetOrderRequest) (*api.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("GetOrder: %v", r))
	return s.next.GetOrder(ctx, r)
}

func (s *TelemetryMiddleware) UpdateOrder(ctx context.Context, o *api.Order) (*api.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("UpdateOrder: %v", o))
	return s.next.UpdateOrder(ctx, o)
}

func (s *TelemetryMiddleware) ValidateOrder(ctx context.Context, items []*api.OrderItem) ([]*api.OrderItem, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("ValidateOrder: %v", items))
	return s.next.ValidateOrder(ctx, items)
}
