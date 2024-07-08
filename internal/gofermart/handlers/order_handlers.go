package handlers

import (
	"encoding/json"
	"gofermart/internal/gofermart/models"
	"gofermart/internal/gofermart/services"
	"gofermart/pkg/helpers"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type OrderHandler struct {
	service *services.OrderService
	logger  *zap.SugaredLogger
}

func NewOrderHandler(service *services.OrderService, logger *zap.SugaredLogger) *OrderHandler {
	return &OrderHandler{service: service, logger: logger}
}

func (h *OrderHandler) LoadOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserIDFromContext(r, h.logger)
	if err != nil {
		h.logger.Errorf("Error getting user id from context: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		h.logger.Error("Error decoding request body", zap.Error(err))
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	order.UserID = userID

	if err := h.service.CreateOrder(r.Context(), &order); err != nil {
		h.logger.Error("Error creating order", zap.Error(err))
		http.Error(w, "Error creating order: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserIDFromContext(r, h.logger)
	if err != nil {
		h.logger.Errorf("Error getting user id from context: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	orders, err := h.service.GetOrdersByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error("Error getting orders", zap.Error(err))
		http.Error(w, "Error getting orders", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, orders)
}
