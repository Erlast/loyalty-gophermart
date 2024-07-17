package mocks

import (
	"context"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) CreateUserTx(ctx context.Context, tx pgx.Tx, user *models.User) error {
	args := m.Called(ctx, tx, user)
	return args.Error(0)
}

func (m *MockUserStore) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	args := m.Called(ctx, login)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserStore) BeginTx(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

type MockBalanceStore struct {
	mock.Mock
}

func (m *MockBalanceStore) GetBalanceByUserID(ctx context.Context, userID int64) (*models.Balance, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.Balance), args.Error(1)
}

func (m *MockBalanceStore) Withdraw(ctx context.Context, withdrawal *models.WithdrawalRequest) error {
	args := m.Called(ctx, withdrawal)
	return args.Error(0)
}

func (m *MockBalanceStore) GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]models.Withdrawal, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Withdrawal), args.Error(1)
}

func (m *MockBalanceStore) CreateBalanceTx(ctx context.Context, tx pgx.Tx, userID int64) error {
	args := m.Called(ctx, tx, userID)
	return args.Error(0)
}

func (m *MockBalanceStore) UpdateBalance(ctx context.Context, userID int64, amount float64) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

func (m *MockBalanceStore) UpdateBalanceTx(ctx context.Context, tx pgx.Tx, userID int64, accrual float64) error {
	args := m.Called(ctx, tx, userID, accrual)
	return args.Error(0)
}
