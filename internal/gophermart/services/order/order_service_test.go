package order

import (
	"context"
	"errors"
	"testing"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/services/order/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
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
