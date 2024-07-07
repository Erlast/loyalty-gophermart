package services

import (
	"context"
	"gofermart/internal/gofermart/models"
	"gofermart/internal/gofermart/storage"
)

type BalanceService struct {
	storage *storage.BalanceStorage
}

func NewBalanceService(storage *storage.BalanceStorage) *BalanceService {
	return &BalanceService{storage: storage}
}

func (s *BalanceService) GetBalanceByUserID(ctx context.Context, userID int64) (*models.Balance, error) {
	return s.storage.GetBalanceByUserID(ctx, userID)
}

func (s *BalanceService) Withdraw(ctx context.Context, withdrawal *models.WithdrawalRequest) error {
	balance, err := s.storage.GetBalanceByUserID(ctx, withdrawal.UserID)
	if err != nil {
		return err
	}

	if balance.CurrentBalance < withdrawal.Amount {
		return errors.New("insufficient balance")
	}

	return s.storage.Withdraw(ctx, withdrawal)
}

func (s *BalanceService) GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]models.Withdrawal, error) {
	return s.storage.GetWithdrawalsByUserID(ctx, userID)
}
