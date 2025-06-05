package main

import (
	"context"
	"fmt"

	"github.com/shimkek/omd-common/api"
	"go.opentelemetry.io/otel/trace"
)

type TelemetryMiddleware struct {
	next PaymentsService
}

func NewTelemetryMiddleware(service PaymentsService) *TelemetryMiddleware {
	return &TelemetryMiddleware{service}
}

func (s *TelemetryMiddleware) CreatePayment(ctx context.Context, order *api.Order) (string, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("CreatePayment: %v", order))
	return s.next.CreatePayment(ctx, order)
}
