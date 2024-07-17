package services

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/storage"
	"github.com/Erlast/loyalty-gophermart.git/pkg/validators"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidOrderNumber  = errors.New("invalid order number format")
)

type BalanceService struct {
	logger  *zap.SugaredLogger
	storage *storage.BalanceStorage
}

func NewBalanceService(
	balanceStorage *storage.BalanceStorage,
	logger *zap.SugaredLogger,
) *BalanceService {
	return &BalanceService{
		logger:  logger,
		storage: balanceStorage,
	}
}

func (s *BalanceService) GetBalanceByUserID(ctx context.Context, userID int64) (*models.Balance, error) {
	balance, err := s.storage.GetBalanceByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting balance by user id %d: %w", userID, err)
	}
	return balance, nil
}

func (s *BalanceService) Withdraw(ctx context.Context, withdrawal *models.WithdrawalRequest) error {
	balance, err := s.storage.GetBalanceByUserID(ctx, withdrawal.UserID)
	if err != nil {
		return fmt.Errorf("error withdraw balance: %w", err)
	}

	if balance.CurrentBalance < withdrawal.Amount {
		return ErrInsufficientBalance
	}

	if !validators.ValidateOrderNumber(withdrawal.Order) {
		return ErrInvalidOrderNumber
	}

	err = s.storage.Withdraw(ctx, withdrawal)
	if err != nil {
		return fmt.Errorf("error withdraw: %w", err)
	}
	return nil
}

func (s *BalanceService) GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]models.Withdrawal, error) {
	withdrawals, err := s.storage.GetWithdrawalsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error balance service getting withdrawals by user id %d: %w", userID, err)
	}
	return withdrawals, nil
}
