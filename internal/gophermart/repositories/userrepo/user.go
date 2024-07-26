package userrepo

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore interface {
	CreateUser(ctx context.Context, user *models.User) error
	CreateUserTx(ctx context.Context, tx pgx.Tx, user *models.User) error
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type UserStorage struct {
	logger *zap.SugaredLogger
	db     *pgxpool.Pool
}

func NewUserStorage(
	db *pgxpool.Pool,
	logger *zap.SugaredLogger,
) *UserStorage {
	return &UserStorage{
		logger: logger,
		db:     db,
	}
}

func tableExists(ctx context.Context, db *pgxpool.Pool, tableName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		);
	`
	err := db.QueryRow(ctx, query, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("could not check if table exists: %w", err)
	}
	return exists, nil
}

func (s *UserStorage) CreateUserTx(ctx context.Context, tx pgx.Tx, user *models.User) error {
	exist, err := tableExists(ctx, s.db, "users")
	if err != nil {
		return fmt.Errorf("error checking if table exists: %w", err)
	}
	if !exist {
		return fmt.Errorf("table users doesn't exists ")
	} else {
		return fmt.Errorf("tablle exists")
	}

	//query := `INSERT INTO users (login, password, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id`
	//err = tx.QueryRow(ctx, query, user.Login, user.Password).Scan(&user.ID)
	//if err != nil {
	//	return fmt.Errorf("could not create user: %w", err)
	//}
	//return nil
}

func (s *UserStorage) CreateUser(ctx context.Context, user *models.User) error {
	// Проверяем наличие таблицы перед вставкой
	exist, err := tableExists(ctx, s.db, "users")
	if err != nil {
		return fmt.Errorf("error checking if table exists: %w", err)
	}
	if !exist {
		return fmt.Errorf("table users doesn't exists ")
	} else {
		return fmt.Errorf("tablle exists")
	}

	// Запрос на вставку данных пользователя в таблицу
	//query := `INSERT INTO users (login, password, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id`
	//err = s.db.QueryRow(ctx, query, user.Login, user.Password).Scan(&user.ID)
	//if err != nil {
	//	return fmt.Errorf("could not create user: %w", err)
	//}
	//return nil
}

func (s *UserStorage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `SELECT id, login, password, created_at, updated_at FROM users WHERE login=$1`
	user := &models.User{}
	err := s.db.QueryRow(ctx, query, login).Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting user by login: %w", err)
	}
	return user, nil
}

func (s *UserStorage) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	return tx, nil
}
