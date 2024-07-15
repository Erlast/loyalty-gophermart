package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userStorage    *storage.UserStorage
	balanceStorage *storage.BalanceStorage
}

func NewUserService(userStorage *storage.UserStorage) *UserService {
	return &UserService{userStorage: userStorage}
}

func (s *UserService) Register(ctx context.Context, user *models.User) error {
	op := "user service register"
	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	user.Password = string(hashedPassword)
	tx, err := s.userStorage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	err = s.userStorage.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	err = s.balanceStorage.CreateBalance(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, credentials models.Credentials) (*models.User, error) {
	op := "user service login"
	user, err := s.userStorage.GetUserByLogin(ctx, credentials.Login)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
