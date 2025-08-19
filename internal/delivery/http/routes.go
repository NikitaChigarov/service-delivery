package http

import (
	"net/http"
	"service-delivery/internal/service"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router, service *service.OrderService) {
	handler := NewOrderHandler(service)

	// API routes
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/orders/{order_uid}", handler.GetOrderByUID).Methods("GET")

	// Web interface
	router.HandleFunc("/", handler.ServeWebInterface)
	router.HandleFunc("/order/{order_uid}", handler.ServeWebInterface)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web"))))
}
