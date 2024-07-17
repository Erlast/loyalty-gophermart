package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	logger         *zap.SugaredLogger
	userStorage    *storage.UserStorage
	balanceStorage *storage.BalanceStorage
}

func NewUserService(
	userStorage *storage.UserStorage,
	balanceStorage *storage.BalanceStorage,
	logger *zap.SugaredLogger,
) *UserService {
	return &UserService{
		logger:         logger,
		userStorage:    userStorage,
		balanceStorage: balanceStorage,
	}
}

func (s *UserService) Register(ctx context.Context, user *models.User) (err error) {
	op := "user service register method"

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	user.Password = string(hashedPassword)

	// Начало транзакции
	tx, err := s.userStorage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			s.logger.Errorf("%s: %w", op, err)
		}
	}(tx, ctx)

	// Создание пользователя в транзакции
	err = s.userStorage.CreateUserTx(ctx, tx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Создание баланса в транзакции
	err = s.balanceStorage.CreateBalanceTx(ctx, tx, user.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
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
