package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/services"
	"github.com/Erlast/loyalty-gophermart.git/pkg/helpers"
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
	h.logger.Infof("User id from context: %v", userID)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Error reading request body", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.logger.Error("Error closing body", zap.Error(err))
		}
	}(r.Body)
	h.logger.Infof("Request body: %v", string(body))
	orderNumber := string(body)

	order := models.Order{
		UserID:     userID,
		Number:     orderNumber,
		Status:     string(models.OrderStatusNew),
		UploadedAt: time.Now(),
	}
	h.logger.Infof("Struct Order: %v", order)

	err = h.service.CreateOrder(ctx, &order)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrOrderAlreadyLoadedBySameUser):
			h.logger.Error("Order number already loaded by this user", zap.Error(err))
			http.Error(w, "", http.StatusOK)
			return
		case errors.Is(err, services.ErrOrderAlreadyLoadedByDifferentUser):
			h.logger.Error("Order number already loaded by a different user", zap.Error(err))
			http.Error(w, "", http.StatusConflict)
			return
		case errors.Is(err, services.ErrInvalidOrderFormat):
			h.logger.Error("Invalid order number format", zap.Error(err))
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		default:
			h.logger.Error("Error creating order", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}

	h.logger.Infof("Created order: %v", order)

	w.WriteHeader(http.StatusAccepted)
}

func (h *OrderHandler) ListOrders(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.logger.Info("List orders called")
	userID, err := helpers.GetUserIDFromContext(r, h.logger)
	if err != nil {
		h.logger.Errorf("Error getting user id from context: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	h.logger.Infof("User id from context: %v", userID)

	orders, err := h.service.GetOrdersByUserID(ctx, userID)
	if err != nil {
		h.logger.Error("Error getting orders", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	h.logger.Infof("List of orders: %v", orders)

	w.Header().Set("Content-Type", "application/json")
	render.JSON(w, r, orders)
}
