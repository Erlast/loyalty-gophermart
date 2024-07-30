package order

import (
	"context"
	"errors"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/services/order/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderStorage := mocks.NewMockOrderStorage(ctrl)
	mockAccrualService := mocks.NewMockAccrualService(ctrl) // Assume no need in this specific test
	logger := zap.NewExample().Sugar()
	os := NewOrderService(mockOrderStorage, nil, mockAccrualService, logger)

	ctx := context.Background()
	testOrder := &models.Order{Number: "12345678903", UserID: 1}

	t.Run("Invalid Order Number", func(t *testing.T) {
		mockOrderStorage.EXPECT().CheckOrder(gomock.Any(), gomock.Any()).Times(0) // Should not check storage
		err := os.CreateOrder(ctx, &models.Order{Number: "invalid", UserID: 1})
		if err == nil || !errors.Is(err, ErrInvalidOrderFormat) {
			t.Errorf("Expected ErrInvalidOrderFormat, got %v", err)
		}
	})

	t.Run("Order Exists", func(t *testing.T) {
		mockOrderStorage.EXPECT().CheckOrder(ctx, "12345678903").
			Return(&models.Order{Number: "12345678903", UserID: 1}, nil)
		err := os.CreateOrder(ctx, testOrder)
		if err == nil || !errors.Is(err, ErrOrderAlreadyLoadedBySameUser) {
			t.Errorf("Expected ErrOrderAlreadyLoadedBySameUser, got %v", err)
		}
	})

	t.Run("Successful Creation", func(t *testing.T) {
		mockOrderStorage.EXPECT().CheckOrder(ctx, "12345678903").Return(nil, pgx.ErrNoRows)
		mockOrderStorage.EXPECT().CreateOrder(ctx, testOrder).Return(nil)
		err := os.CreateOrder(ctx, testOrder)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestUpdateOrderStatuses_ErrorGettingOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockOrderStorage := mocks.NewMockOrderStorage(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	mockAccrualService := mocks.NewMockAccrualService(ctrl)
	mockBalanceStorage := mocks.NewMockBalanceStorage(ctrl)
	logger := zap.NewExample().Sugar()
	os := NewOrderService(mockOrderStorage, mockBalanceStorage, mockAccrualService, logger)

	t.Run("ErrorGetOrdersByStatus", func(t *testing.T) {
		var orders []models.Order
		mockOrderStorage.EXPECT().GetOrdersByStatus(ctx, models.OrderStatusNew, models.OrderStatusProcessing).
			Return(orders, errors.New("test error"))
		err := os.UpdateOrderStatuses(ctx)
		assert.EqualError(t, err, "error getting orders: test error")
	})

	t.Run("ErrorGetBeginTx", func(t *testing.T) {
		var mockOrders []models.Order
		mockOrderStorage.EXPECT().GetOrdersByStatus(ctx, models.OrderStatusNew, models.OrderStatusProcessing).
			Return(mockOrders, nil)
		mockOrderStorage.EXPECT().BeginTx(ctx).Return(mockTx, errors.New("test error2"))
		err := os.UpdateOrderStatuses(ctx)
		assert.EqualError(t, err, "error starting transaction: test error2")
	})

	t.Run("SuccessCase", func(t *testing.T) {
		timeForTest := time.Now()
		uploadedAt1 := timeForTest.Add(-1 * time.Hour)
		uploadedAt2 := timeForTest.Add(-1 * time.Second)
		mockOrders := []models.Order{
			{
				ID:         1,
				Number:     "123452670",
				UserID:     100,
				Status:     "NEW",
				UploadedAt: uploadedAt1,
			},
			{
				ID:         2,
				Number:     "12345678111",
				UserID:     100,
				Status:     "PROCESSING",
				UploadedAt: uploadedAt2,
			},
		}
		mockOrderStorage.EXPECT().GetOrdersByStatus(ctx, models.OrderStatusNew, models.OrderStatusProcessing).
			Return(mockOrders, nil)
		mockOrderStorage.EXPECT().BeginTx(ctx).Return(mockTx, nil)
		accrual1 := float32(300)
		accrual2 := float32(500)
		mockAccrualService.EXPECT().GetAccrualInfo("123452670").
			Return(
				&models.AccrualResponse{
					Order:   "123452670",
					Status:  "PROCESSED",
					Accrual: accrual1,
				}, nil)

		mockAccrualService.EXPECT().GetAccrualInfo("12345678111").
			Return(
				&models.AccrualResponse{
					Order:   "12345678111",
					Status:  "PROCESSED",
					Accrual: accrual2,
				}, nil)

		mockOrders[0].Status = "PROCESSED"
		mockOrders[0].Accrual = &accrual1
		mockOrderStorage.EXPECT().UpdateOrderTx(ctx, mockTx, &mockOrders[0]).Return(nil)

		mockOrders[1].Status = "PROCESSED"
		mockOrders[1].Accrual = &accrual2
		mockOrderStorage.EXPECT().UpdateOrderTx(ctx, mockTx, &mockOrders[1]).Return(nil)

		mockBalanceStorage.EXPECT().UpdateBalanceTx(ctx, mockTx, mockOrders[0].UserID, *mockOrders[0].Accrual).Return(nil)
		mockBalanceStorage.EXPECT().UpdateBalanceTx(ctx, mockTx, mockOrders[1].UserID, *mockOrders[1].Accrual).Return(nil)
		mockTx.EXPECT().Commit(ctx).Return(nil)

		err := os.UpdateOrderStatuses(ctx)
		require.NoError(t, err)
	})
}
