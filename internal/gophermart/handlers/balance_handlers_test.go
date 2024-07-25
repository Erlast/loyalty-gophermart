package handlers

import (
	"context"
	"encoding/json"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/repositories/balance"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/services"
	"github.com/Erlast/loyalty-gophermart.git/pkg/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBalanceHandler_GetBalance(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	balanceMockStorage := new(balance.MockBalanceStore)
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
