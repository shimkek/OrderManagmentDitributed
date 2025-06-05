package main

import (
	"context"
	"time"

	"github.com/shimkek/omd-common/api"
	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	next OrderService
}

func NewLoggingMiddleware(next OrderService) *LoggingMiddleware {
	return &LoggingMiddleware{next}
}

func (s *LoggingMiddleware) CreateOrder(ctx context.Context, p *api.CreateOrderRequest) (*api.Order, error) {
	start := time.Now()

	defer func() {
		zap.L().Info("CreateOrder", zap.Duration("took", time.Since(start)))
	}()
	return s.next.CreateOrder(ctx, p)
}

func (s *LoggingMiddleware) GetOrder(ctx context.Context, r *api.GetOrderRequest) (*api.Order, error) {
	start := time.Now()

	defer func() {
		zap.L().Info("GetOrder", zap.Duration("took", time.Since(start)))
	}()

	return s.next.GetOrder(ctx, r)
}

func (s *LoggingMiddleware) UpdateOrder(ctx context.Context, o *api.Order) (*api.Order, error) {
	start := time.Now()

	defer func() {
		zap.L().Info("UpdateOrder", zap.Duration("took", time.Since(start)))
	}()

	return s.next.UpdateOrder(ctx, o)
}

func (s *LoggingMiddleware) ValidateOrder(ctx context.Context, items []*api.OrderItem) ([]*api.OrderItem, error) {
	start := time.Now()

	defer func() {
		zap.L().Info("ValidateOrder", zap.Duration("took", time.Since(start)))
	}()
	return s.next.ValidateOrder(ctx, items)
}
