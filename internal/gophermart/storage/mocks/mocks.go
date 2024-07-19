package mocks

import (
	"context"
	"fmt"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) CreateUserTx(ctx context.Context, tx pgx.Tx, user *models.User) error {
	args := m.Called(ctx, tx, user)
	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("error MockUserStore CreateUserTx: %w", err)
	}
	return nil
}

func (m *MockUserStore) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	args := m.Called(ctx, login)
	user, ok := args.Get(0).(*models.User)
	if !ok {
		return nil, fmt.Errorf("error in GetUserByLogin: expected *models.User, got %T", args.Get(0))
	}
	err := args.Error(1)
	if err != nil {
		err = fmt.Errorf("error MockUserStore GetUserByLogin: %w", err)
	}
	return user, err
}

func (m *MockUserStore) BeginTx(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	tx, ok := args.Get(0).(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf("error in BeginTx: expected pgx.Tx, got %T", args.Get(0))
	}
	err := args.Error(1)
	if err != nil {
		err = fmt.Errorf("error MockUserStore BeginTx: %w", err)
	}
	return tx, err
}

type MockBalanceStore struct {
	mock.Mock
}

func (m *MockBalanceStore) GetBalanceByUserID(ctx context.Context, userID int64) (*models.Balance, error) {
	args := m.Called(ctx, userID)
	balance, ok := args.Get(0).(*models.Balance)
	if !ok {
		return nil, fmt.Errorf("error in GetBalanceByUserID: expected *models.Balance, got %T", args.Get(0))
	}
	err := args.Error(1)
	if err != nil {
		err = fmt.Errorf("error MockBalanceStore GetBalanceByUserID: %w", err)
	}
	return balance, err
}

func (m *MockBalanceStore) Withdraw(ctx context.Context, withdrawal *models.WithdrawalRequest) error {
	args := m.Called(ctx, withdrawal)
	err := args.Error(1)
	if err != nil {
		err = fmt.Errorf("error MockBalanceStore Withdraw: %w", err)
	}
	return err
}

func (m *MockBalanceStore) GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]models.Withdrawal, error) {
	args := m.Called(ctx, userID)
	withdrawals, ok := args.Get(0).([]models.Withdrawal)
	if !ok {
		return nil, fmt.Errorf("error in GetWithdrawalsByUserID: expected []models.Withdrawal, got %T", args.Get(0))
	}
	err := args.Error(1)
	if err != nil {
		err = fmt.Errorf("error MockBalanceStore GetWithdrawalsByUserID: %w", err)
	}
	return withdrawals, err
}

func (m *MockBalanceStore) CreateBalanceTx(ctx context.Context, tx pgx.Tx, userID int64) error {
	args := m.Called(ctx, tx, userID)
	err := args.Error(0)
	if err != nil {
		err = fmt.Errorf("error MockBalanceStore CreateBalanceTx: %w", err)
	}
	return err
}

func (m *MockBalanceStore) UpdateBalance(ctx context.Context, userID int64, amount float32) error {
	args := m.Called(ctx, userID, amount)
	err := args.Error(0)
	if err != nil {
		err = fmt.Errorf("error MockBalanceStore UpdateBalance: %w", err)
	}
	return err
}

func (m *MockBalanceStore) UpdateBalanceTx(ctx context.Context, tx pgx.Tx, userID int64, accrual float32) error {
	args := m.Called(ctx, tx, userID, accrual)
	err := args.Error(0)
	if err != nil {
		err = fmt.Errorf("error MockBalanceStore UpdateBalanceTx: %w", err)
	}
	return err
}
