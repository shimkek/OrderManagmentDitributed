package inmem

import "github.com/shimkek/omd-common/api"

type inmemProcessor struct{}

func NewInmemProcessor() *inmemProcessor {
	return &inmemProcessor{}
}

func (p *inmemProcessor) CreatePaymentLink(*api.Order) (string, error) {
	// In-memory implementation just returns a dummy link
	return "http://example.com/payment-link", nil
}
