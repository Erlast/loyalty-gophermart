package storage

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
)

const errString = "err %w"

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
	args := m.Called(ctx, orderNumber)
	data, ok := args.Get(0).(*models.Order)
	if !ok {
		return nil, fmt.Errorf(errString, args.Error(1))
	}
	return data, nil
}

func (m *MockStorage) SaveOrderItems(ctx context.Context, items models.OrderItem) error {
	args := m.Called(ctx, items)
	if args.Error(0) != nil {
		return fmt.Errorf(errString, args.Error(0))
	}
	return nil
}

func (m *MockStorage) SaveGoods(ctx context.Context, goods models.Goods) error {
	args := m.Called(ctx, goods)
	if args.Error(0) != nil {
		return fmt.Errorf(errString, args.Error(0))
	}
	return nil
}

func (m *MockStorage) GetRegisteredOrders(ctx context.Context) ([]int64, error) {
	args := m.Called(ctx)
	data, ok := args.Get(0).([]int64)
	if !ok {
		return nil, fmt.Errorf(errString, args.Error(1))
	}
	return data, nil
}

func (m *MockStorage) FetchRewardRules(ctx context.Context) ([]models.Goods, error) {
	args := m.Called(ctx)
	data, ok := args.Get(0).([]models.Goods)
	if !ok {
		return nil, fmt.Errorf(errString, args.Error(1))
	}
	return data, nil
}

func (m *MockStorage) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	args := m.Called(ctx, orderID, status)
	if args.Error(0) != nil {
		return fmt.Errorf(errString, args.Error(0))
	}
	return nil
}

func (m *MockStorage) FetchProducts(ctx context.Context, orderID int64) ([]models.Items, error) {
	args := m.Called(ctx, orderID)
	data, ok := args.Get(0).([]models.Items)
	if !ok {
		return nil, fmt.Errorf(errString, args.Error(1))
	}
	return data, nil
}

func (m *MockStorage) SaveOrderPoints(ctx context.Context, orderID int64, points []float32) error {
	args := m.Called(ctx, orderID, points)
	if args.Error(0) != nil {
		return fmt.Errorf(errString, args.Error(0))
	}
	return nil
}

func (m *MockStorage) GetProcessedOrders(ctx context.Context) ([]models.Order, error) {
	args := m.Called(ctx)
	data, ok := args.Get(0).([]models.Order)
	if !ok {
		return nil, fmt.Errorf(errString, args.Error(1))
	}
	return data, nil
}
