package services

import (
	"context"
	"errors"
	"fmt"
	"gofermart/internal/gofermart/models"

	"gofermart/internal/gofermart/storage"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	storage *storage.UserStorage
}

func NewUserService(userStorage *storage.UserStorage) *UserService {
	return &UserService{storage: userStorage}
}

func (s *UserService) Register(ctx context.Context, user *models.User) error {
	op := "user service register"
	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	user.Password = string(hashedPassword)
	err = s.storage.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (s *UserService) Login(ctx context.Context, credentials models.Credentials) (*models.User, error) {
	op := "user service login"
	user, err := s.storage.GetUserByLogin(ctx, credentials.Login)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
