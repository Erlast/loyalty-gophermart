package storage

import (
	"context"
	"gofermart/internal/gofermart/models"

	"github.com/jackc/pgx/v4/pgxpool"
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

func (s *OrderStorage) GetOrdersByUserID(ctx context.Context, userID int64) ([]models.Order, error) {
	query := `SELECT id, user_id, number, status, accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at`
	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
