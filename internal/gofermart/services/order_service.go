package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/Erlast/loyalty-gophermart.git/internal/gofermart/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/gofermart/storage"
	"github.com/Erlast/loyalty-gophermart.git/pkg/validators"

	"time"
)

var (
	ErrOrderAlreadyLoadedBySameUser      = errors.New("order number already loaded by this user")
	ErrOrderAlreadyLoadedByDifferentUser = errors.New("order number already loaded by a different user")
	ErrInvalidOrderFormat                = errors.New("invalid order number format")
)

type OrderService struct {
	storage        *storage.OrderStorage
	accrualService *AccrualService
}

func NewOrderService(orderStorage *storage.OrderStorage, accrualService *AccrualService) *OrderService {
	return &OrderService{storage: orderStorage, accrualService: accrualService}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	existingOrder, err := s.storage.GetOrder(ctx, order.Number)
	if err != nil {
		return fmt.Errorf("error getting existing order: %w", err)
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

	if err = s.storage.CreateOrder(ctx, order); err != nil {
		return fmt.Errorf("error creating order: %w", err)
	}
	return nil
}

func (s *OrderService) GetOrdersByUserID(ctx context.Context, userID int64) ([]models.Order, error) {
	order, err := s.storage.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting orders: %w", err)
	}
	return order, nil
}

func (s *OrderService) UpdateOrderStatuses(ctx context.Context) error {
	orders, err := s.storage.GetOrdersByStatus(ctx, models.OrderStatusNew, models.OrderStatusProcessing)
	if err != nil {
		return fmt.Errorf("error getting orders: %w", err)
	}

	for _, order := range orders {
		accrualInfo, err := s.accrualService.GetAccrualInfo(order.Number)
		if err != nil {
			continue
		}

		if accrualInfo == nil {
			continue
		}

		order.Status = accrualInfo.Status
		order.Accrual = accrualInfo.Accrual
		order.UploadedAt = time.Now()

		if err := s.storage.UpdateOrder(ctx, &order); err != nil {
			return fmt.Errorf("error updating order: %w", err)
		}
	}

	return nil
}

func (s *OrderService) GetOrderAccrualInfo(ctx context.Context, orderNumber string) (*models.AccrualResponse, error) {
	accrualInfo, err := s.accrualService.GetAccrualInfo(orderNumber)
	if err != nil {
		return nil, fmt.Errorf("error getting accrual info: %w", err)
	}

	if accrualInfo == nil {
		return nil, models.ErrOrderNotFound
	}

	return accrualInfo, nil
}
