package handlers

import (
	"context"
	"errors"
	"gofermart/internal/gofermart/models"
	"gofermart/internal/gofermart/services"
	"gofermart/pkg/helpers"
	"gofermart/pkg/validators"
	"io"
	"net/http"
	"time"

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

func (h *OrderHandler) LoadOrder(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserIDFromContext(r, h.logger)
	if err != nil {
		h.logger.Errorf("Error getting user id from context: %v", err)
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Error reading request body", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	orderNumber := string(body)

	if !validators.ValidateOrderNumber(orderNumber) { // Assuming you have a function to validate the order number
		h.logger.Error("Invalid order number format", zap.String("orderNumber", orderNumber))
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}

	order := models.Order{
		UserID:     userID,
		Number:     orderNumber,
		Status:     string(models.OrderStatusNew),
		UploadedAt: time.Now(),
	}

	err = h.service.CreateOrder(ctx, &order)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrOrderAlreadyLoadedBySameUser):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Order number already loaded by this user")) //nolint:errcheck
		case errors.Is(err, services.ErrOrderAlreadyLoadedByDifferentUser):
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Order number already loaded by a different user")) //nolint:errcheck
		case errors.Is(err, services.ErrInvalidOrderFormat):
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Invalid order number format")) //nolint:errcheck
		default:
			h.logger.Error("Error creating order", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *OrderHandler) ListOrders(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserIDFromContext(r, h.logger)
	if err != nil {
		h.logger.Errorf("Error getting user id from context: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	orders, err := h.service.GetOrdersByUserID(ctx, userID)
	if err != nil {
		h.logger.Error("Error getting orders", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, orders)
}
