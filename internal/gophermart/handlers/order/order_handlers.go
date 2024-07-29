package order

import (
	"context"
	"errors"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/middleware"
	"io"
	"net/http"
	"time"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/services/order"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type OrderHandler struct {
	service *order.OrderService
	logger  *zap.SugaredLogger
}

func NewOrderHandler(service *order.OrderService, logger *zap.SugaredLogger) *OrderHandler {
	return &OrderHandler{service: service, logger: logger}
}

func (h *OrderHandler) LoadOrder(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
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

	orderStruct := models.Order{
		UserID:     userID,
		Number:     orderNumber,
		Status:     string(models.OrderStatusNew),
		UploadedAt: time.Now(),
	}
	h.logger.Infof("Struct Order: %v", orderStruct)

	err = h.service.CreateOrder(ctx, &orderStruct)
	if err != nil {
		switch {
		case errors.Is(err, order.ErrOrderAlreadyLoadedBySameUser):
			h.logger.Error("Order number already loaded by this user", zap.Error(err))
			http.Error(w, "", http.StatusOK)
			return
		case errors.Is(err, order.ErrOrderAlreadyLoadedByDifferentUser):
			h.logger.Errorf("Order number already loaded by a different user: %v", orderStruct)
			http.Error(w, "", http.StatusConflict)
			return
		case errors.Is(err, order.ErrInvalidOrderFormat):
			h.logger.Error("Invalid order number format", zap.Error(err))
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		default:
			h.logger.Error("Error creating order", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}

	h.logger.Infof("Created order: %v", orderStruct)

	w.WriteHeader(http.StatusAccepted)
}

func (h *OrderHandler) ListOrders(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.logger.Info("List orders called")
	userID, err := middleware.GetUserIDFromContext(r.Context())
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
