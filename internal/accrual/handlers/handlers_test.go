package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
	"github.com/Erlast/loyalty-gophermart.git/pkg/zaplog"
)

func TestGetAccrualOrderHandler(t *testing.T) {
	store := &storage.MockStorage{}
	newLogger := zaplog.InitLogger()

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

	GetAccrualOrderHandler(context.Background(), res, req, store, newLogger)

	assert.Equal(t, http.StatusOK, res.Code)

	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var order models.Order
	err := json.NewDecoder(res.Body).Decode(&order)
	assert.NoError(t, err)
	assert.Equal(t, "1234567812345670", order.UUID)
	assert.Equal(t, "PROCESSED", order.Status)
}

func TestGetAccrualOrderHandler_NotFound(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	store := &storage.MockStorage{}

	store.On("GetByOrderNumber", mock.Anything, "invalid-order-number").Return(nil, errors.New("not found"))

	req, err := http.NewRequest(http.MethodGet, "/orders/invalid-order-number", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	r := chi.NewRouter()

	r.Get("/orders/{number}", func(w http.ResponseWriter, r *http.Request) {
		GetAccrualOrderHandler(r.Context(), w, r, store, logger)
	})

	r.ServeHTTP(rec, req)

	resp := rec.Result()
	err = resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestPostAccrualOrderHandler(t *testing.T) {
	store := &storage.MockStorage{}
	newLogger := zaplog.InitLogger()

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

	PostAccrualOrderHandler(context.Background(), res, req, store, newLogger)

	assert.Equal(t, http.StatusAccepted, res.Code)

	store.AssertCalled(t, "SaveOrderItems", req.Context(), orderItem)
}

func TestPostAccrualOrderHandler_Conflict(t *testing.T) {
	store := &storage.MockStorage{}
	newLogger := zaplog.InitLogger()

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

	PostAccrualOrderHandler(context.Background(), res, req, store, newLogger)

	assert.Equal(t, http.StatusConflict, res.Code)

	store.AssertCalled(t, "SaveOrderItems", req.Context(), orderItem)
}

func TestPostAccrualOrderHandler_InternalServerError(t *testing.T) {
	store := &storage.MockStorage{}
	newLogger := zaplog.InitLogger()

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

	PostAccrualOrderHandler(context.Background(), res, req, store, newLogger)

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	store.AssertCalled(t, "SaveOrderItems", req.Context(), orderItem)
}

func TestPostAccrualGoodsHandler(t *testing.T) {
	store := &storage.MockStorage{}
	newLogger := zaplog.InitLogger()

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

	PostAccrualGoodsHandler(context.Background(), res, req, store, newLogger)

	assert.Equal(t, http.StatusOK, res.Code)

	store.AssertCalled(t, "SaveGoods", req.Context(), goods)
}

func TestPostAccrualGoodsHandler_Conflict(t *testing.T) {
	store := &storage.MockStorage{}
	newLogger := zaplog.InitLogger()

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

	PostAccrualGoodsHandler(context.Background(), res, req, store, newLogger)

	assert.Equal(t, http.StatusConflict, res.Code)

	store.AssertCalled(t, "SaveGoods", req.Context(), goods)
}

func TestPostAccrualGoodsHandler_InternalServerError(t *testing.T) {
	store := &storage.MockStorage{}
	newLogger := zaplog.InitLogger()

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

	PostAccrualGoodsHandler(context.Background(), res, req, store, newLogger)

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	store.AssertCalled(t, "SaveGoods", req.Context(), goods)
}
