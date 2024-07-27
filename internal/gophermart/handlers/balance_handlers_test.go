package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/repositories/balancerepo"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/services"
	"github.com/Erlast/loyalty-gophermart.git/pkg/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

func TestBalanceHandler_GetBalance_StatusOK(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	balanceMockStorage := new(balancerepo.MockBalanceStore)
	balanceModel := models.Balance{
		UserID:         1,
		CurrentBalance: 100,
		TotalWithdrawn: 900,
	}
	balanceMockStorage.On("GetBalanceByUserID", mock.Anything, int64(1)).
		Return(&balanceModel, nil)

	balanceService := services.NewBalanceService(balanceMockStorage, logger)
	balanceHandler := NewBalanceHandler(balanceService, logger)

	ctx := context.Background()
	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", http.NoBody)
	rr := httptest.NewRecorder()

	frCtx := helpers.NewMockFormContext()
	frCtx.On("GetUserID", req, logger).Return(int64(1), nil)

	balanceHandler.GetBalance(ctx, rr, req, frCtx)

	assert.Equal(t, http.StatusOK, rr.Code)

	var respBalance models.Balance
	err := json.Unmarshal(rr.Body.Bytes(), &respBalance)
	if err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	assert.Equal(t, balanceModel.CurrentBalance, respBalance.CurrentBalance)

	balanceMockStorage.AssertExpectations(t)
	frCtx.AssertExpectations(t)
}

func TestBalanceHandler_GetBalance_InternalServerError1(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	balanceMockStorage := new(balancerepo.MockBalanceStore)
	balanceModel := models.Balance{}

	balanceMockStorage.On("GetBalanceByUserID", mock.Anything, int64(1)).
		Return(&balanceModel, errors.New("not can get balance by user"))

	balanceService := services.NewBalanceService(balanceMockStorage, logger)
	balanceHandler := NewBalanceHandler(balanceService, logger)

	ctx := context.Background()
	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", http.NoBody)
	rr := httptest.NewRecorder()

	frCtx := helpers.NewMockFormContext()
	frCtx.On("GetUserID", req, logger).Return(int64(1), nil)

	balanceHandler.GetBalance(ctx, rr, req, frCtx)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	balanceMockStorage.AssertExpectations(t)
	frCtx.AssertExpectations(t)
}

func TestBalanceHandler_GetBalance_InternalServerError2(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	balanceMockStorage := new(balancerepo.MockBalanceStore)

	balanceService := services.NewBalanceService(balanceMockStorage, logger)
	balanceHandler := NewBalanceHandler(balanceService, logger)

	ctx := context.Background()
	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", http.NoBody)
	rr := httptest.NewRecorder()

	frCtx := helpers.NewMockFormContext()
	frCtx.On("GetUserID", req, logger).Return(int64(1), errors.New("user not found"))

	balanceHandler.GetBalance(ctx, rr, req, frCtx)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	frCtx.AssertExpectations(t)
}

func TestBalanceHandler_Withdrawals_StatusOK(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t).Sugar()
	balanceMockStorage := new(balancerepo.MockBalanceStore)

	withdrawals := []models.Withdrawal{
		{
			ProcessedAt: time.Now().Add(-1 * time.Hour),
			Order:       "100",
			Amount:      100,
			UserID:      2,
		},
		{
			ProcessedAt: time.Now().Add(-24 * time.Hour),
			Order:       "100",
			Amount:      100,
			UserID:      2,
		},
	}

	balanceMockStorage.On("GetWithdrawalsByUserID", ctx, int64(2)).
		Return(withdrawals, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", http.NoBody)
	// Установка заголовка Content-Type
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	frCtx := helpers.NewMockFormContext()
	frCtx.On("GetUserID", req, logger).Return(int64(2), nil)

	balanceService := services.NewBalanceService(balanceMockStorage, logger)
	balanceHandler := NewBalanceHandler(balanceService, logger)
	balanceHandler.Withdrawals(ctx, rr, req, frCtx)

	jsonString := rr.Body.String()

	var outputWithdrawals []models.Withdrawal
	// Unmarshal JSON into the slice of Order structures
	err := json.Unmarshal([]byte(jsonString), &outputWithdrawals)
	require.NoError(t, err)
	ow := len(outputWithdrawals)

	assert.Equal(t, 2, ow)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestBalanceHandler_Withdrawals_InternalServerError(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t).Sugar()
	balanceMockStorage := new(balancerepo.MockBalanceStore)

	withdrawals := []models.Order{
		{
			Number: "123",
			Status: "NEW",
			UserID: 2,
		},
		{
			Number: "456",
			Status: "NEW",
			UserID: 2,
		},
	}

	balanceMockStorage.On("GetWithdrawalsByUserID", ctx, int64(2)).
		Return(withdrawals, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", http.NoBody)
	// Установка заголовка Content-Type
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	frCtx := helpers.NewMockFormContext()
	frCtx.On("GetUserID", req, logger).Return(int64(2), nil)

	balanceService := services.NewBalanceService(balanceMockStorage, logger)
	balanceHandler := NewBalanceHandler(balanceService, logger)
	balanceHandler.Withdrawals(ctx, rr, req, frCtx)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestBalanceHandler_Withdrawals_NoContent(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t).Sugar()
	balanceMockStorage := new(balancerepo.MockBalanceStore)

	var withdrawals = make([]models.Withdrawal, 0)
	balanceMockStorage.On("GetWithdrawalsByUserID", ctx, int64(2)).
		Return(withdrawals, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", http.NoBody)
	// Установка заголовка Content-Type
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	frCtx := helpers.NewMockFormContext()
	frCtx.On("GetUserID", req, logger).Return(int64(2), nil)

	balanceService := services.NewBalanceService(balanceMockStorage, logger)
	balanceHandler := NewBalanceHandler(balanceService, logger)
	balanceHandler.Withdrawals(ctx, rr, req, frCtx)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestBalanceHandler_Withdrawals_(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t).Sugar()
	balanceMockStorage := new(balancerepo.MockBalanceStore)

	var withdrawals = make([]models.Withdrawal, 0)
	balanceMockStorage.On("GetWithdrawalsByUserID", ctx, int64(2)).
		Return(withdrawals, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", http.NoBody)
	// Установка заголовка Content-Type
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	frCtx := helpers.NewMockFormContext()
	frCtx.On("GetUserID", req, logger).Return(int64(2), errors.New("user not found"))

	balanceService := services.NewBalanceService(balanceMockStorage, logger)
	balanceHandler := NewBalanceHandler(balanceService, logger)
	balanceHandler.Withdrawals(ctx, rr, req, frCtx)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
