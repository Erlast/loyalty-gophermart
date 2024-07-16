package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/storage"
	"github.com/Erlast/loyalty-gophermart.git/pkg/zaplog"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

const (
	errorDescription    = "%s: %w"
	errorDescriptionLog = "%s: %v"
)

type UserService struct {
	userStorage    *storage.UserStorage
	balanceStorage *storage.BalanceStorage
}

func NewUserService(userStorage *storage.UserStorage, balanceStorage *storage.BalanceStorage) *UserService {
	return &UserService{
		userStorage:    userStorage,
		balanceStorage: balanceStorage,
	}
}

func (s *UserService) Register(ctx context.Context, user *models.User) error {
	op := "user service register"
	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf(errorDescription, op, err)
	}

	user.Password = string(hashedPassword)

	err = s.userStorage.CreateUser(ctx, user)
	if err != nil {
		zaplog.Logger.Errorf("%s: %w", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.balanceStorage.CreateBalance(ctx, user.ID)
	if err != nil {
		zaplog.Logger.Errorf(errorDescriptionLog, op, err)
		return fmt.Errorf(errorDescription, op, err)
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

// IsDuplicateError checks if the error is a duplicate entry error.
func IsDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // Unique violation error code in PostgreSQL
	}
	return false
}
