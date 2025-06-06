package main

import (
	"fmt"
	"log"
	"net/http"

	common "github.com/shimkek/omd-common"
	"github.com/shimkek/omd-common/api"
	"github.com/shimkek/omd-gateway/gateway"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type handler struct {
	gateway gateway.OrdersGateway
}

func NewHandler(gateway gateway.OrdersGateway) *handler {
	return &handler{gateway: gateway}
}
func (h *handler) registerRoutes(mux *http.ServeMux) {
	//static routes serving
	mux.Handle("/", http.FileServer(http.Dir("public")))

	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)
	mux.HandleFunc("GET /api/customers/{customerID}/orders/{orderID}", h.HandleGetOrder)
}
func (h *handler) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	orderID := r.PathValue("orderID")
	if customerID == "" || orderID == "" {
		http.Error(w, "customerID and orderID are required", http.StatusBadRequest)
		return
	}
	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	order, err := h.gateway.GetOrder(ctx, orderID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		log.Print("failed to get order:", err)
		common.WriteError(w, http.StatusInternalServerError, "failed to get order")
		return
	}
	if err := common.WriteJson(w, http.StatusOK, &order); err != nil {
		common.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
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
	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.gateway.CreateOrder(ctx, &api.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		log.Print("failed to create order:", err)
		common.WriteError(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	res := &CreateOrderRequest{
		Order:         order,
		RedirectToURL: fmt.Sprintf("http://localhost:8080/success.html?customerID=%s&orderID=%s", order.CustomerID, order.OrderID),
	}
	if err := common.WriteJson(w, http.StatusCreated, res); err != nil {
		common.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
}

func validateItems(items []*api.OrderItem) error {
	if len(items) == 0 {
		return fmt.Errorf("items cannot be empty")
	}
	for _, item := range items {
		if item.ProductID == "" {
			return fmt.Errorf("productID is required")
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("quantity must be greater than 0")
		}
	}
	return nil
}
