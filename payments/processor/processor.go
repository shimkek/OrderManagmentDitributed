package processor

import "github.com/shimkek/omd-common/api"

type PaymentProcessor interface {
	CreatePaymentLink(*api.Order) (string, error)
}
