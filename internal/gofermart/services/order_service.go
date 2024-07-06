package services

import (
	"context"
	"errors"
	"gofermart/internal/gofermart/models"
	"gofermart/internal/gofermart/storage"
	"gofermart/pkg/validators"
	"regexp"
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

// validateOrderNumber проверяет, что номер заказа состоит только из цифр и имеет длину от 10 до 20 символов.
func validateOrderNumber(orderNumber string) bool {
	re := regexp.MustCompile(`^\d{10,20}$`)
	return re.MatchString(orderNumber)
}
