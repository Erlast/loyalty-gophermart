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
	op := "order handler method load order called"

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
			w.WriteHeader(http.StatusOK)
			_, err = w.Write([]byte("Order number already loaded by this user"))
			if err != nil {
				h.logger.Errorf("%v:,%v", op, err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		case errors.Is(err, services.ErrOrderAlreadyLoadedByDifferentUser):
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte("Order number already loaded by a different user"))
			if err != nil {
				h.logger.Errorf("%v:,%v", op, err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		case errors.Is(err, services.ErrInvalidOrderFormat):
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, err = w.Write([]byte(InvalidOrderFormatMsg))
			if err != nil {
				h.logger.Errorf("%v:,%v", op, err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		default:
			h.logger.Error("Error creating order", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		return
	}

	h.logger.Infof("Created order: %v", order)

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
