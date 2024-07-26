package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/repositories/balancerepo"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/repositories/userrepo"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	logger         *zap.SugaredLogger
	userStorage    userrepo.UserStore // Используем интерфейс UserStore
	balanceStorage balancerepo.BalanceStore
}

func NewUserService(
	userStorage userrepo.UserStore, // Используем интерфейс UserStore
	balanceStorage balancerepo.BalanceStore,
	logger *zap.SugaredLogger,
) *UserService {
	return &UserService{
		logger:         logger,
		userStorage:    userStorage,
		balanceStorage: balanceStorage,
	}
}

func (s *UserService) Register(ctx context.Context, user *models.User) (err error) {
	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}

	user.Password = string(hashedPassword)

	// Создание пользователя в транзакции
	err = s.userStorage.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}

	// Создание баланса в транзакции
	err = s.balanceStorage.CreateBalance(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error creating balance: %w", err)
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, credentials models.Credentials) (*models.User, error) {
	user, err := s.userStorage.GetUserByLogin(ctx, credentials.Login)
	if err != nil {
		return nil, fmt.Errorf("could not get user by login: %w", err)
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
