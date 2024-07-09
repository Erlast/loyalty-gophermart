package services

import (
	"context"
	"errors"
	"gofermart/internal/gofermart/models"
	"gofermart/internal/gofermart/storage"
	"gofermart/pkg/validators"
)

var (
	ErrOrderAlreadyLoadedBySameUser      = errors.New("order number already loaded by this user")
	ErrOrderAlreadyLoadedByDifferentUser = errors.New("order number already loaded by a different user")
	ErrInvalidOrderFormat                = errors.New("invalid order number format")
)

type OrderService struct {
	storage *storage.OrderStorage
}

func NewOrderService(storage *storage.OrderStorage) *OrderService {
	return &OrderService{storage: storage}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	existingOrder, err := s.storage.GetOrder(ctx, order.Number)
	if err != nil {
		return err
	}

	if existingOrder != nil {
		if existingOrder.UserID == order.UserID {
			return ErrOrderAlreadyLoadedBySameUser
		}
		return ErrOrderAlreadyLoadedByDifferentUser
	}

	if !validators.ValidateOrderNumber(order.Number) {
		return ErrInvalidOrderFormat
	}

	return s.storage.CreateOrder(ctx, order)
}

func (s *OrderService) GetOrdersByUserID(ctx context.Context, userID int64) ([]models.Order, error) {
	return s.storage.GetOrdersByUserID(ctx, userID)
}
