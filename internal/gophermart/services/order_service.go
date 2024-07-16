package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Erlast/loyalty-gophermart.git/pkg/zaplog"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/storage"
	"github.com/Erlast/loyalty-gophermart.git/pkg/validators"

	"time"
)

var (
	ErrOrderAlreadyLoadedBySameUser      = errors.New("order number already loaded by this user")
	ErrOrderAlreadyLoadedByDifferentUser = errors.New("order number already loaded by a different user")
	ErrInvalidOrderFormat                = errors.New("invalid order number format")
)

type OrderService struct {
	orderStorage   *storage.OrderStorage
	balanceStorage *storage.BalanceStorage
	accrualService *AccrualService
}

func NewOrderService(
	orderStorage *storage.OrderStorage,
	balanceStorage *storage.BalanceStorage,
	accrualService *AccrualService,
) *OrderService {
	return &OrderService{
		orderStorage:   orderStorage,
		balanceStorage: balanceStorage,
		accrualService: accrualService,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	existOrder, err := s.orderStorage.CheckOrder(ctx, order.Number)

	if existOrder == true {
		return ErrOrderAlreadyLoadedByDifferentUser
	}

	if !validators.ValidateOrderNumber(order.Number) {
		return ErrInvalidOrderFormat
	}

	if err = s.orderStorage.CreateOrder(ctx, order); err != nil {
		return fmt.Errorf("error creating order: %w", err)
	}
	return nil
}

func (s *OrderService) GetOrdersByUserID(ctx context.Context, userID int64) ([]models.Order, error) {
	order, err := s.orderStorage.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting orders: %w", err)
	}
	return order, nil
}

func (s *OrderService) UpdateOrderStatuses(ctx context.Context) error {
	orders, err := s.orderStorage.GetOrdersByStatus(ctx, models.OrderStatusNew, models.OrderStatusProcessing)
	zaplog.Logger.Infof("orders: %v", orders)
	if err != nil {
		return fmt.Errorf("error getting orders: %w", err)
	}

	for _, order := range orders {
		zaplog.Logger.Infof("GetAccrualInfo number: %v", order.Number)
		accrualInfo, err := s.accrualService.GetAccrualInfo(order.Number)
		if err != nil {
			continue
		}

		if accrualInfo == nil {
			continue
		}

		order.Status = accrualInfo.Status
		order.Accrual = &accrualInfo.Accrual
		order.UploadedAt = time.Now()
		zaplog.Logger.Infof("Order struct for update: %v", order)

		if err := s.orderStorage.UpdateOrder(ctx, &order); err != nil {
			return fmt.Errorf("error updating order: %w", err)
		}

		if order.Status == string(models.OrderStatusProcessed) {
			err := s.balanceStorage.UpdateBalance(ctx, order.UserID, *order.Accrual)
			if err != nil {
				zaplog.Logger.Infof("Error updating balance: %v", err)
				return fmt.Errorf("error updating balance: %w", err)
			}
		}
	}

	return nil
}
