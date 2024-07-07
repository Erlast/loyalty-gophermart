package services

import (
	"context"
	"errors"
	"gofermart/internal/gofermart/models"
	"gofermart/internal/gofermart/storage"
	"gofermart/pkg/validators"
)

type OrderService struct {
	storage *storage.OrderStorage
}

func NewOrderService(storage *storage.OrderStorage) *OrderService {
	return &OrderService{storage: storage}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	if !validators.ValidateOrderNumber(order.Number) {
		return errors.New("invalid order number")
	}

	return s.storage.CreateOrder(ctx, order)
}

func (s *OrderService) GetOrdersByUserID(ctx context.Context, userID int64) ([]models.Order, error) {
	return s.storage.GetOrdersByUserID(ctx, userID)
}
