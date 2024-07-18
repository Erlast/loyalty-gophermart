package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/storage"
	"github.com/Erlast/loyalty-gophermart.git/pkg/validators"
)

var (
	ErrOrderAlreadyLoadedBySameUser      = errors.New("order number already loaded by this user")
	ErrOrderAlreadyLoadedByDifferentUser = errors.New("order number already loaded by a different user")
	ErrInvalidOrderFormat                = errors.New("invalid order number format")
)

type OrderService struct {
	logger         *zap.SugaredLogger
	orderStorage   *storage.OrderStorage
	balanceStorage *storage.BalanceStorage
	accrualService *AccrualService
}

func NewOrderService(
	orderStorage *storage.OrderStorage,
	balanceStorage *storage.BalanceStorage,
	accrualService *AccrualService,
	logger *zap.SugaredLogger,
) *OrderService {
	return &OrderService{
		logger:         logger,
		orderStorage:   orderStorage,
		balanceStorage: balanceStorage,
		accrualService: accrualService,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	existOrder, err := s.orderStorage.CheckOrder(ctx, order.Number)
	if err != nil {
		return fmt.Errorf("error checking order: %w", err)
	}

	if existOrder {
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
	orders, err := s.orderStorage.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting orders: %w", err)
	}
	return orders, nil
}

func (s *OrderService) UpdateOrderStatuses(ctx context.Context) error {
	orders, err := s.orderStorage.GetOrdersByStatus(ctx, models.OrderStatusNew, models.OrderStatusProcessing)
	s.logger.Infof("orders: %v", orders)
	if err != nil {
		return fmt.Errorf("error getting orders: %w", err)
	}

	tx, err := s.orderStorage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			s.logger.Infof("Error rolling back transaction: %v", err)
		}
	}(tx, ctx)

	for _, order := range orders {
		s.logger.Infof("GetAccrualInfo number: %v", order.Number)
		accrualInfo, err := s.accrualService.GetAccrualInfo(order.Number)
		if err != nil {
			s.logger.Errorf("error getting accrualInfo: %v", err)
			continue
		}

		if accrualInfo == nil {
			continue
		}

		order.Status = accrualInfo.Status
		order.Accrual = &accrualInfo.Accrual
		s.logger.Infof("Order struct for update: %v", order)

		if err := s.orderStorage.UpdateOrderTx(ctx, tx, &order); err != nil {
			return fmt.Errorf("error updating order: %w", err)
		}

		if order.Status == string(models.OrderStatusProcessed) {
			err := s.balanceStorage.UpdateBalanceTx(ctx, tx, order.UserID, *order.Accrual)
			if err != nil {
				s.logger.Infof("Error updating balance: %v", err)
				return fmt.Errorf("error updating balance: %w", err)
			}
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (s *OrderService) UpdateOrderStatusesOld(ctx context.Context) error {
	orders, err := s.orderStorage.GetOrdersByStatus(ctx, models.OrderStatusNew, models.OrderStatusProcessing)
	s.logger.Infof("orders: %v", orders)
	if err != nil {
		return fmt.Errorf("error getting orders: %w", err)
	}

	for _, order := range orders {
		s.logger.Infof("GetAccrualInfo number: %v", order.Number)
		accrualInfo, err := s.accrualService.GetAccrualInfo(order.Number)
		if err != nil {
			s.logger.Errorf("error getting accrualInfo: %v", err)
			continue
		}

		if accrualInfo == nil {
			continue
		}

		order.Status = accrualInfo.Status
		order.Accrual = &accrualInfo.Accrual
		s.logger.Infof("Order struct for update: %v", order)

		if err := s.orderStorage.UpdateOrder(ctx, &order); err != nil {
			return fmt.Errorf("error updating order: %w", err)
		}

		if order.Status == string(models.OrderStatusProcessed) {
			err := s.balanceStorage.UpdateBalance(ctx, order.UserID, *order.Accrual)
			if err != nil {
				s.logger.Infof("Error updating balance: %v", err)
				return fmt.Errorf("error updating balance: %w", err)
			}
		}
	}

	return nil
}
