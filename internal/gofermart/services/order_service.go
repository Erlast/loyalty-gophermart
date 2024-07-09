package services

import (
	"context"
	"errors"
	"gofermart/internal/gofermart/models"
	"gofermart/internal/gofermart/storage"
	"gofermart/pkg/validators"
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

func NewOrderService(storage *storage.OrderStorage, accrualService *AccrualService) *OrderService {
	return &OrderService{storage: storage, accrualService: accrualService}
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

func (s *OrderService) UpdateOrderStatuses(ctx context.Context) error {
	orders, err := s.storage.GetOrdersByStatus(ctx, models.OrderStatusNew, models.OrderStatusProcessing)
	if err != nil {
		return err
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
			return err
		}
	}

	return nil
}

func (s *OrderService) GetOrderAccrualInfo(ctx context.Context, orderNumber string) (*models.AccrualResponse, error) {
	accrualInfo, err := s.accrualService.GetAccrualInfo(orderNumber)
	if err != nil {
		return nil, err
	}

	if accrualInfo == nil {
		return nil, models.ErrOrderNotFound
	}

	return accrualInfo, nil
}
