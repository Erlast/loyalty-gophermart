package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
	"github.com/Erlast/loyalty-gophermart.git/pkg/zaplog"
)

func TestGetAccrualOrderHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := NewMockStorage(ctrl)
	newLogger := zaplog.InitLogger()

	req := httptest.NewRequest(http.MethodGet, "/orders/1234567812345670", http.NoBody)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{"number"},
			Values: []string{"1234567812345670"},
		},
	}))

	res := httptest.NewRecorder()

	store.EXPECT().GetByOrderNumber(req.Context(), "1234567812345670").Return(&models.Order{
		UUID:   "1234567812345670",
		Status: "PROCESSED",
	}, nil)

	GetAccrualOrderHandler(res, req, store, newLogger)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := NewMockStorage(ctrl)

	store.EXPECT().GetByOrderNumber(gomock.Any(), "invalid-order-number").Return(nil, errors.New("not found"))

	req, err := http.NewRequest(http.MethodGet, "/orders/invalid-order-number", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	r := chi.NewRouter()

	r.Get("/orders/{number}", func(w http.ResponseWriter, r *http.Request) {
		GetAccrualOrderHandler(w, r, store, logger)
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := NewMockStorage(ctrl)
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

	store.EXPECT().SaveOrderItems(req.Context(), orderItem).Return(nil)

	PostAccrualOrderHandler(res, req, store, newLogger)

	assert.Equal(t, http.StatusAccepted, res.Code)
}

func TestPostAccrualOrderHandler_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := NewMockStorage(ctrl)
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

	store.EXPECT().SaveOrderItems(req.Context(), orderItem).Return(&helpers.ConflictError{})

	PostAccrualOrderHandler(res, req, store, newLogger)

	assert.Equal(t, http.StatusConflict, res.Code)
}

func TestPostAccrualOrderHandler_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := NewMockStorage(ctrl)
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

	store.EXPECT().SaveOrderItems(req.Context(), orderItem).Return(errors.New("internal error"))

	PostAccrualOrderHandler(res, req, store, newLogger)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestPostAccrualGoodsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := NewMockStorage(ctrl)
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

	store.EXPECT().SaveGoods(req.Context(), goods).Return(nil)

	PostAccrualGoodsHandler(res, req, store, newLogger)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestPostAccrualGoodsHandler_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := NewMockStorage(ctrl)
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

	store.EXPECT().SaveGoods(req.Context(), goods).Return(&helpers.ConflictError{})

	PostAccrualGoodsHandler(res, req, store, newLogger)

	assert.Equal(t, http.StatusConflict, res.Code)
}

func TestPostAccrualGoodsHandler_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := NewMockStorage(ctrl)
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

	store.EXPECT().SaveGoods(req.Context(), goods).Return(errors.New("internal error"))

	PostAccrualGoodsHandler(res, req, store, newLogger)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}
