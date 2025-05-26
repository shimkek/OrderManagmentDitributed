package main

import (
	"context"

	"github.com/shimkek/omd-common/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *api.Order) (string, error)
}
