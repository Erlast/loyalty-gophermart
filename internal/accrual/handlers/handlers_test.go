package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/Erlast/loyalty-gophermart.git/pkg/zaplog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetAccrualOrderHandler(t *testing.T) {
	store := &storage.MockStorage{}
	zaplog.InitLogger()

	req := httptest.NewRequest(http.MethodGet, "/orders/1234567812345670", http.NoBody)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{"number"},
			Values: []string{"1234567812345670"},
		},
	}))

	res := httptest.NewRecorder()

	store.On("GetByOrderNumber", req.Context(), "1234567812345670").Return(&models.Order{
		UUID:   "1234567812345670",
		Status: "PROCESSED",
	}, nil)

	GetAccrualOrderHandler(context.Background(), res, req, store)

	assert.Equal(t, http.StatusOK, res.Code)

	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var order models.Order
	err := json.NewDecoder(res.Body).Decode(&order)
	assert.NoError(t, err)
	assert.Equal(t, "1234567812345670", order.UUID)
	assert.Equal(t, "PROCESSED", order.Status)
}

func TestPostAccrualOrderHandler(t *testing.T) {
	store := &storage.MockStorage{}
	zaplog.InitLogger()

	var goods []models.Items

	goods = append(goods, models.Items{
		Description: "Чайник Bork",
		Price:       700,
	})
	orderItem := models.OrderItem{
		UUID:  "1234567812345670",
		Goods: goods,
	}
	body, _ := json.Marshal(orderItem)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	store.On("SaveOrderItems", req.Context(), orderItem).Return(nil)

	PostAccrualOrderHandler(context.Background(), res, req, store)

	assert.Equal(t, http.StatusAccepted, res.Code)

	store.AssertCalled(t, "SaveOrderItems", req.Context(), orderItem)
}

func TestPostAccrualOrderHandler_Conflict(t *testing.T) {
	store := &storage.MockStorage{}
	zaplog.InitLogger()

	var goods []models.Items

	goods = append(goods, models.Items{
		Description: "Чайник Bork",
		Price:       700,
	})
	orderItem := models.OrderItem{
		UUID:  "1234567812345670",
		Goods: goods,
	}
	body, _ := json.Marshal(orderItem)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	store.On("SaveOrderItems", req.Context(), orderItem).Return(&helpers.ConflictError{})

	PostAccrualOrderHandler(context.Background(), res, req, store)

	assert.Equal(t, http.StatusConflict, res.Code)

	store.AssertCalled(t, "SaveOrderItems", req.Context(), orderItem)
}

func TestPostAccrualOrderHandler_InternalServerError(t *testing.T) {
	store := &storage.MockStorage{}
	zaplog.InitLogger()

	var goods []models.Items

	goods = append(goods, models.Items{
		Description: "Чайник Bork",
		Price:       700,
	})
	orderItem := models.OrderItem{
		UUID:  "1234567812345670",
		Goods: goods,
	}
	body, _ := json.Marshal(orderItem)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	store.On("SaveOrderItems", req.Context(), orderItem).Return(errors.New("internal error"))

	PostAccrualOrderHandler(context.Background(), res, req, store)

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	store.AssertCalled(t, "SaveOrderItems", req.Context(), orderItem)
}

func TestPostAccrualGoodsHandler(t *testing.T) {
	store := &storage.MockStorage{}
	zaplog.InitLogger()

	goods := models.Goods{
		Match:      "somebrand",
		Reward:     10,
		RewardType: "%",
	}
	body, _ := json.Marshal(goods)
	req := httptest.NewRequest(http.MethodPost, "/goods", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	store.On("SaveGoods", req.Context(), goods).Return(nil)

	PostAccrualGoodsHandler(context.Background(), res, req, store)

	assert.Equal(t, http.StatusOK, res.Code)

	store.AssertCalled(t, "SaveGoods", req.Context(), goods)
}

func TestPostAccrualGoodsHandler_Conflict(t *testing.T) {
	store := &storage.MockStorage{}
	zaplog.InitLogger()

	goods := models.Goods{
		Match:      "somebrand",
		Reward:     10,
		RewardType: "%",
	}
	body, _ := json.Marshal(goods)
	req := httptest.NewRequest(http.MethodPost, "/goods", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	store.On("SaveGoods", req.Context(), goods).Return(&helpers.ConflictError{})

	PostAccrualGoodsHandler(context.Background(), res, req, store)

	assert.Equal(t, http.StatusConflict, res.Code)

	store.AssertCalled(t, "SaveGoods", req.Context(), goods)
}

func TestPostAccrualGoodsHandler_InternalServerError(t *testing.T) {
	store := &storage.MockStorage{}
	zaplog.InitLogger()

	goods := models.Goods{
		Match:      "somebrand",
		Reward:     10,
		RewardType: "%",
	}
	body, _ := json.Marshal(goods)
	req := httptest.NewRequest(http.MethodPost, "/goods", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	store.On("SaveGoods", req.Context(), goods).Return(errors.New("internal error"))

	PostAccrualGoodsHandler(context.Background(), res, req, store)

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	store.AssertCalled(t, "SaveGoods", req.Context(), goods)
}
