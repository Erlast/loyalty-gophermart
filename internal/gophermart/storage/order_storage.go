package storage

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderStorage struct {
	logger *zap.SugaredLogger
	db     *pgxpool.Pool
}

func NewOrderStorage(
	db *pgxpool.Pool,
	logger *zap.SugaredLogger,
) *OrderStorage {
	return &OrderStorage{
		logger: logger,
		db:     db,
	}
}

func (s *OrderStorage) GetOrder(ctx context.Context, number string) (*models.Order, error) {
	var order models.Order
	query := "SELECT user_id, number FROM orders WHERE number=$1"
	row := s.db.QueryRow(ctx, query, number)
	err := row.Scan(&order.UserID, &order.Number)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("order not found: %w", err)
		}
		return nil, fmt.Errorf("error order storage getting order: %w", err)
	}
	return &order, nil
}

func (s *OrderStorage) CheckOrder(ctx context.Context, number string) (bool, error) {
	query := "SELECT 1 FROM orders WHERE number=$1"
	row := s.db.QueryRow(ctx, query, number)
	var exists int
	err := row.Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error checking order existence: %w", err)
	}
	return true, nil
}

func (s *OrderStorage) CreateOrder(ctx context.Context, order *models.Order) error {
	query := "INSERT INTO orders (user_id, number, status, uploaded_at) VALUES ($1, $2, $3, $4)"
	s.logger.Infof("inserting new order: %v,%v,%v,%v", order.UserID, order.Number, order.Status, order.UploadedAt)
	_, err := s.db.Exec(ctx, query, order.UserID, order.Number, order.Status, order.UploadedAt)
	if err != nil {
		return fmt.Errorf("error order storage create order: %w", err)
	}
	return nil
}

func (s *OrderStorage) GetOrdersByUserID(ctx context.Context, userID int64) ([]models.Order, error) {
	query := `SELECT number, status, accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at`
	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error order storage get orders by user: %w", err)
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
			return nil, fmt.Errorf("error scan order get orders by user: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error read rows get orders by user: %w", err)
	}

	return orders, nil
}

func (s *OrderStorage) GetOrdersByStatus(ctx context.Context, statuses ...models.OrderStatus) ([]models.Order, error) {
	var orders []models.Order
	query := "SELECT user_id, number, status, accrual, uploaded_at FROM orders WHERE status = ANY($1)"
	rows, err := s.db.Query(ctx, query, statuses)
	if err != nil {
		return nil, fmt.Errorf("error order storage get orders by status: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.UserID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, fmt.Errorf("error scan order get orders by status: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error read rows get orders by status: %w", err)
	}

	return orders, nil
}

func (s *OrderStorage) UpdateOrder(ctx context.Context, order *models.Order) error {
	query := "UPDATE orders SET status=$1, accrual=ROUND($2, 2) WHERE number=$3"
	_, err := s.db.Exec(ctx, query, order.Status, order.Accrual, order.Number)
	if err != nil {
		return fmt.Errorf("error order storage update order: %w", err)
	}
	return nil
}

func (s *OrderStorage) UpdateOrderTx(ctx context.Context, tx pgx.Tx, order *models.Order) error {
	query := "UPDATE orders SET status=$1, accrual=ROUND($2, 2) WHERE number=$3"
	_, err := tx.Exec(ctx, query, order.Status, order.Accrual, order.Number)
	if err != nil {
		return fmt.Errorf("error updating order: %v", err)
	}
	return nil
}

func (s *OrderStorage) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error begin transaction from order storage: %w", err)
	}
	return tx, nil
}
