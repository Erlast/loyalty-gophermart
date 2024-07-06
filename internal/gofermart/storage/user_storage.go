package storage

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"gofermart/internal/gofermart/models"
)

type UserStorage struct {
	db *pgxpool.Pool
}

func NewUserStorage(db *pgxpool.Pool) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (login, password, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id`
	return s.db.QueryRow(ctx, query, user.Login, user.Password).Scan(&user.ID)
}

func (s *UserStorage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `SELECT id, login, password, created_at, updated_at FROM users WHERE login=$1`
	user := &models.User{}
	err := s.db.QueryRow(ctx, query, login).Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
