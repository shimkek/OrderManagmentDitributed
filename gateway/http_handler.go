package main

import "net/http"

type handler struct {
	//gateway instance
}

func NewHandler() *handler {
	return &handler{}
}
func (h *handler) registerRoutes(mux *http.ServeMux) error {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)

	return nil
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {

}
