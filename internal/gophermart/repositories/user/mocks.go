package user

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
