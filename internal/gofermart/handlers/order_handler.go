package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"gofermart/internal/gofermart/models"
	"gofermart/internal/gofermart/services"
	"net/http"
)

type OrderHandler struct {
	Service *services.OrderService
	Logger  *zap.SugaredLogger
}

func NewOrderHandler(service *services.OrderService, logger *zap.SugaredLogger) *OrderHandler {
	return &OrderHandler{Service: service, Logger: logger}
}

func (h *OrderHandler) LoadOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		h.Logger.Error("Error decoding request body", zap.Error(err))
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	order.UserID = userID

	if err := h.Service.CreateOrder(r.Context(), &order); err != nil {
		h.Logger.Error("Error creating order", zap.Error(err))
		http.Error(w, "Error creating order: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
