package services

import (
	"context"
	"errors"
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
	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.storage.CreateUser(ctx, user)
}

func (s *UserService) Login(ctx context.Context, credentials models.Credentials) (*models.User, error) {
	user, err := s.storage.GetUserByLogin(ctx, credentials.Login)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
