package services

import (
	"context"
	"errors"
	"fmt"
	"gofermart/internal/gofermart/models"
	"gofermart/internal/gofermart/storage"
	"gofermart/pkg/validators"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidOrderNumber  = errors.New("invalid order number format")
)

type BalanceService struct {
	storage *storage.BalanceStorage
}

func NewBalanceService(balanceStorage *storage.BalanceStorage) *BalanceService {
	return &BalanceService{storage: balanceStorage}
}

func (s *BalanceService) GetBalanceByUserID(ctx context.Context, userID int64) (*models.Balance, error) {
	return s.storage.GetBalanceByUserID(ctx, userID)
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

	return s.storage.Withdraw(ctx, withdrawal)
}

func (s *BalanceService) GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]models.Withdrawal, error) {
	return s.storage.GetWithdrawalsByUserID(ctx, userID)
}
