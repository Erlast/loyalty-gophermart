package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/repositories/balancerepo"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/repositories/userrepo"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type MockTx struct {
	mock.Mock
	pgx.Tx
}

func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("error MockTx Commit: %w", err)
	}
	return nil
}

func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("error MockTx Rollback: %w", err)
	}
	return nil
}

func TestUserService_Register(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop().Sugar()

	mockUserStore := new(userrepo.MockUserStore)
	mockBalanceStore := new(balancerepo.MockBalanceStore)
	mockTx := new(MockTx)

	userService := NewUserService(mockUserStore, mockBalanceStore, logger)

	user := &models.User{
		Login:    "testuser",
		Password: "password",
	}

	mockUserStore.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockUserStore.On("CreateUserTx", mock.Anything, mockTx, user).Return(nil)
	mockBalanceStore.On("CreateBalanceTx", mock.Anything, mockTx, user.ID).Return(nil)
	mockTx.On("Commit", mock.Anything).Return(nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)

	err := userService.Register(ctx, user)
	assert.NoError(t, err)

	mockUserStore.AssertExpectations(t)
	mockBalanceStore.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUserService_Login(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop().Sugar()

	mockUserStore := new(userrepo.MockUserStore)
	mockBalanceStore := new(balancerepo.MockBalanceStore)

	userService := NewUserService(mockUserStore, mockBalanceStore, logger)

	credentials := models.Credentials{
		Login:    "testuser",
		Password: "password",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	user := &models.User{
		Login:    credentials.Login,
		Password: string(hashedPassword),
	}

	mockUserStore.On("GetUserByLogin", mock.Anything, credentials.Login).Return(user, nil)

	result, err := userService.Login(ctx, credentials)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.Login, result.Login)

	mockUserStore.AssertExpectations(t)
}
