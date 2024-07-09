package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"gofermart/internal/gofermart/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderStorage struct {
	db *pgxpool.Pool
}

func NewOrderStorage(db *pgxpool.Pool) *OrderStorage {
	return &OrderStorage{db: db}
}

func (s *OrderStorage) GetOrder(ctx context.Context, number string) (*models.Order, error) {
	var order models.Order
	query := "SELECT user_id, number FROM orders WHERE number=$1"
	row := s.db.QueryRow(ctx, query, number)
	err := row.Scan(&order.UserID, &order.Number)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (s *OrderStorage) CreateOrder(ctx context.Context, order *models.Order) error {
	query := "INSERT INTO orders (user_id, number, status, uploaded_at) VALUES ($1, $2, $3, $4)"
	_, err := s.db.Exec(ctx, query, order.UserID, order.Number, order.Status, order.UploadedAt)
	return err
}

func (s *OrderStorage) GetOrdersByUserID(ctx context.Context, userID int64) ([]models.Order, error) {
	query := `SELECT number, status, accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at`
	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *OrderStorage) GetOrdersByStatus(ctx context.Context, statuses ...models.OrderStatus) ([]models.Order, error) {
	var orders []models.Order
	query := "SELECT user_id, number, status, accrual, created_at, updated_at FROM orders WHERE status = ANY($1)"
	rows, err := s.db.Query(ctx, query, statuses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.UserID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *OrderStorage) UpdateOrder(ctx context.Context, order *models.Order) error {
	query := "UPDATE orders SET status=$1, accrual=$2, updated_at=$3 WHERE number=$4"
	_, err := s.db.Exec(ctx, query, order.Status, order.Accrual, order.UploadedAt, order.Number)
	return err
}
