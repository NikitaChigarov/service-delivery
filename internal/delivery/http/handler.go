package http

import (
	"encoding/json"
	"net/http"
	"service-delivery/internal/service"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) GetOrderByUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["order_uid"]

	order, err := h.service.GetOrderByUID(r.Context(), orderUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if order == nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *OrderHandler) ServeWebInterface(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}
