package main

import (
	"context"
	"fmt"

	"github.com/shimkek/omd-common/api"
	"go.opentelemetry.io/otel/trace"
)

type TelemetryMiddleware struct {
	next StockService
}

func NewTelemetryMiddleware(next StockService) *TelemetryMiddleware {
	return &TelemetryMiddleware{next: next}
}

func (s *TelemetryMiddleware) CheckIfItemsAreInStock(ctx context.Context, p []*api.OrderItem) (bool, []*api.OrderItem, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("CheckIfItemAreInStock: %v", p))

	return s.next.CheckIfItemsAreInStock(ctx, p)

}

func (s *TelemetryMiddleware) GetItems(ctx context.Context, ids []string) ([]*api.OrderItem, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("GetItems: %v", ids))

	return s.next.GetItems(ctx, ids)
}
