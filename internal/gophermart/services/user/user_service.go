package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	CreateUserTx(ctx context.Context, tx pgx.Tx, user *models.User) error
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type BalanceStore interface {
	CreateBalanceTx(ctx context.Context, tx pgx.Tx, userID int64) error
}

type UserService struct {
	logger         *zap.SugaredLogger
	userStorage    UserStore
	balanceStorage BalanceStore
}

func NewUserService(
	userStorage UserStore,
	balanceStorage BalanceStore,
	logger *zap.SugaredLogger,
) *UserService {
	return &UserService{
		logger:         logger,
		userStorage:    userStorage,
		balanceStorage: balanceStorage,
	}
}

func (s *UserService) Register(ctx context.Context, user *models.User) error {
	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}

	user.Password = string(hashedPassword)

	// Начало транзакции
	tx, err := s.userStorage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	// Определяем defer сразу после успешного начала транзакции
	defer func() {
		// Проверяем, была ли уже ошибка или ошибка при коммите
		if p := recover(); p != nil || err != nil {
			// В случае ошибки пытаемся откатить транзакцию
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				s.logger.Errorw("Failed to rollback transaction", "err", rbErr)
			}
		}
	}()

	// Создание пользователя в транзакции
	err = s.userStorage.CreateUserTx(ctx, tx, user)
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}

	// Создание баланса в транзакции
	err = s.balanceStorage.CreateBalanceTx(ctx, tx, user.ID)
	if err != nil {
		return fmt.Errorf("error creating balance: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("could not commit: %w", err)
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, credentials models.Credentials) (*models.User, error) {
	user, err := s.userStorage.GetUserByLogin(ctx, credentials.Login)
	if err != nil {
		return nil, fmt.Errorf("could not get user by login: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return nil, fmt.Errorf("could not compare password: %w", err)
	}

	return user, nil
}

// IsDuplicateError checks if the error is a duplicate entry error.
func (s *UserService) IsDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // Unique violation error code in PostgreSQL
	}
	return false
}
