package main

import (
	"net/http"

	common "github.com/shimkek/omd-common"
	"github.com/shimkek/omd-common/api"
)

type handler struct {
	//gateway instance
	orderServiceClient api.OrderServiceClient
}

func NewHandler(orderServiceClient api.OrderServiceClient) *handler {
	return &handler{orderServiceClient: orderServiceClient}
}
func (h *handler) registerRoutes(mux *http.ServeMux) error {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)

	return nil
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	if customerID == "" {
		http.Error(w, "customerID is required", http.StatusBadRequest)
		return
	}

	var items []*api.OrderItem
	if err := common.ReadJson(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	h.orderServiceClient.CreateOrder(r.Context(), &api.CreateOrderRequest{
		CustomerId: customerID,
		Items:      items,
	})
}
