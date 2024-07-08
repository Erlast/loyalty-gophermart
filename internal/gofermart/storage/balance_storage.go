package storage

import (
	"context"
	"gofermart/internal/gofermart/models"
	"gofermart/pkg/zaplog"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v4/pgxpool"
)

type BalanceStorage struct {
	db *pgxpool.Pool
}

func NewBalanceStorage(db *pgxpool.Pool) *BalanceStorage {
	return &BalanceStorage{db: db}
}

func (s *BalanceStorage) GetBalanceByUserID(ctx context.Context, userID int64) (*models.Balance, error) {
	query := `SELECT user_id, current_balance, total_withdrawn FROM balances WHERE user_id = $1`
	row := s.db.QueryRow(ctx, query, userID)

	var balance models.Balance
	if err := row.Scan(&balance.UserID, &balance.CurrentBalance, &balance.TotalWithdrawn); err != nil {
		return nil, err
	}

	return &balance, nil
}

func (s *BalanceStorage) Withdraw(ctx context.Context, withdrawal *models.WithdrawalRequest) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			zaplog.Logger.Error("Rollback failed", zap.Error(err))
		}
	}(tx, ctx)

	updateBalanceQuery := `
        UPDATE balances
        SET current_balance = current_balance - $1, total_withdrawn = total_withdrawn + $1
        WHERE user_id = $2`
	_, err = tx.Exec(ctx, updateBalanceQuery, withdrawal.Amount, withdrawal.UserID)
	if err != nil {
		return err
	}

	insertWithdrawalQuery := `
        INSERT INTO withdrawals (user_id, amount, processed_at)
        VALUES ($1, $2, NOW())`
	_, err = tx.Exec(ctx, insertWithdrawalQuery, withdrawal.UserID, withdrawal.Amount)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *BalanceStorage) GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]models.Withdrawal, error) {
	query := `SELECT user_id, order_number, amount, processed_at FROM withdrawals WHERE user_id = $1 ORDER BY processed_at`
	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []models.Withdrawal
	for rows.Next() {
		var withdrawal models.Withdrawal
		if err := rows.Scan(&withdrawal.UserID, &withdrawal.Order, &withdrawal.Amount, &withdrawal.ProcessedAt); err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, nil
}
