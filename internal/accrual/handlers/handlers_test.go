package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetAccrualOrderHandler(t *testing.T) {
	store := &storage.MockStorage{}
	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	req := httptest.NewRequest(http.MethodGet, "/orders/123", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{"number"},
			Values: []string{"123"},
		},
	}))

	res := httptest.NewRecorder()

	store.On("GetByOrderNumber", req.Context(), "123").Return(&models.Order{
		UUID:   "123",
		Status: "PROCESSED",
	}, nil)

	GetAccrualOrderHandler(context.Background(), res, req, store, logger)

	assert.Equal(t, http.StatusOK, res.Code)

	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var order models.Order
	err := json.NewDecoder(res.Body).Decode(&order)
	assert.NoError(t, err)
	assert.Equal(t, "123", order.UUID)
	assert.Equal(t, "PROCESSED", order.Status)
}

func TestPostAccrualOrderHandler(t *testing.T) {
	store := &storage.MockStorage{}
	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	var goods []models.Items

	goods = append(goods, models.Items{
		Description: "Чайник Bork",
		Price:       700,
	})
	orderItem := models.OrderItem{
		UUID:  "1245",
		Goods: goods,
	}
	body, _ := json.Marshal(orderItem)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	store.On("SaveOrderItems", req.Context(), orderItem).Return(nil)

	PostAccrualOrderHandler(context.Background(), res, req, store, logger)

	assert.Equal(t, http.StatusAccepted, res.Code)

	store.AssertCalled(t, "SaveOrderItems", req.Context(), orderItem)
}

func TestPostAccrualOrderHandler_Conflict(t *testing.T) {
	store := &storage.MockStorage{}
	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	var goods []models.Items

	goods = append(goods, models.Items{
		Description: "Чайник Bork",
		Price:       700,
	})
	orderItem := models.OrderItem{
		UUID:  "1245",
		Goods: goods,
	}
	body, _ := json.Marshal(orderItem)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	store.On("SaveOrderItems", req.Context(), orderItem).Return(&helpers.ConflictError{})

	PostAccrualOrderHandler(context.Background(), res, req, store, logger)

	assert.Equal(t, http.StatusConflict, res.Code)

	store.AssertCalled(t, "SaveOrderItems", req.Context(), orderItem)
}

func TestPostAccrualOrderHandler_InternalServerError(t *testing.T) {
	store := &storage.MockStorage{}
	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	var goods []models.Items

	goods = append(goods, models.Items{
		Description: "Чайник Bork",
		Price:       700,
	})
	orderItem := models.OrderItem{
		UUID:  "1245",
		Goods: goods,
	}
	body, _ := json.Marshal(orderItem)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	store.On("SaveOrderItems", req.Context(), orderItem).Return(errors.New("internal error"))

	PostAccrualOrderHandler(context.Background(), res, req, store, logger)

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	store.AssertCalled(t, "SaveOrderItems", req.Context(), orderItem)
}
