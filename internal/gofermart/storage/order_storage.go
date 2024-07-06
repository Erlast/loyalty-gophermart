package storage

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"gofermart/internal/gofermart/models"
)

type OrderStorage struct {
	db *pgxpool.Pool
}

func NewOrderStorage(db *pgxpool.Pool) *OrderStorage {
	return &OrderStorage{db: db}
}

func (s *OrderStorage) CreateOrder(ctx context.Context, order *models.Order) error {
	query := `INSERT INTO orders (user_id, number, status, uploaded_at) VALUES ($1, $2, 'NEW', NOW())`
	_, err := s.db.Exec(ctx, query, order.UserID, order.Number)
	return err
}
