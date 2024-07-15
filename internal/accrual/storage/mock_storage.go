package storage

import (
	"context"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
	"github.com/stretchr/testify/mock"
)

// MockStorage - мок хранилища
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
	args := m.Called(ctx, orderNumber)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockStorage) SaveOrderItems(ctx context.Context, items models.OrderItem) error {
	//args := m.Called(ctx, orderNumber)
	return nil
}

func (m *MockStorage) SaveGoods(ctx context.Context, goods models.Goods) error {
	//args := m.Called(ctx, orderNumber)
	return nil
}
