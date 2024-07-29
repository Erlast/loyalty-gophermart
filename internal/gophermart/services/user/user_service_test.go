package user

import (
	"context"
	"testing"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := NewMockUserStore(ctrl)
	mockBalanceStore := NewMockBalanceStore(ctrl)
	mockTx := NewMockTx(ctrl)
	logger := zap.NewExample().Sugar()
	userService := NewUserService(mockUserStore, mockBalanceStore, logger)

	ctx := context.Background()
	user := &models.User{
		Login:    "testuser",
		Password: "password123",
	}

	mockUserStore.EXPECT().
		BeginTx(ctx).
		Return(mockTx, nil).
		Times(1)

	mockTx.EXPECT().
		Commit(ctx).
		Return(nil).
		Times(1)

	mockTx.EXPECT().
		Rollback(ctx).
		Return(nil).
		AnyTimes()

	mockUserStore.EXPECT().
		CreateUserTx(ctx, mockTx, gomock.Any()).
		DoAndReturn(func(_ context.Context, _ pgx.Tx, u *models.User) error {
			err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("password123"))
			if err != nil {
				t.Errorf("Password was not hashed correctly")
			}
			return nil
		}).
		Times(1)

	mockBalanceStore.EXPECT().
		CreateBalanceTx(ctx, mockTx, gomock.Any()).
		Return(nil).
		Times(1)

	if err := userService.Register(ctx, user); err != nil {
		t.Errorf("Register failed: %v", err)
	}
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := NewMockUserStore(ctrl)
	logger := zap.NewExample().Sugar()
	userService := NewUserService(mockUserStore, nil, logger)

	ctx := context.Background()
	credentials := models.Credentials{
		Login:    "testuser",
		Password: "password123",
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUserStore.EXPECT().
		GetUserByLogin(ctx, "testuser").
		Return(&models.User{Login: "testuser", Password: string(hashedPassword)}, nil).
		Times(1)

	if _, err := userService.Login(ctx, credentials); err != nil {
		t.Errorf("Login failed: %v", err)
	}
}
