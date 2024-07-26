package balancerepo

import (
	"context"
	"fmt"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BalanceStore interface {
	GetBalanceByUserID(ctx context.Context, userID int64) (*models.Balance, error)
	Withdraw(ctx context.Context, withdrawal *models.WithdrawalRequest) error
	GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]models.Withdrawal, error)
	CreateBalance(ctx context.Context, userID int64) error
	CreateBalanceTx(ctx context.Context, tx pgx.Tx, userID int64) error
	UpdateBalance(ctx context.Context, userID int64, amount float32) error
	UpdateBalanceTx(ctx context.Context, tx pgx.Tx, userID int64, accrual float32) error
}

type BalanceStorage struct {
	logger *zap.SugaredLogger
	db     *pgxpool.Pool
}

func NewBalanceStorage(
	db *pgxpool.Pool,
	logger *zap.SugaredLogger,
) *BalanceStorage {
	return &BalanceStorage{
		logger: logger,
		db:     db,
	}
}

func (s *BalanceStorage) GetBalanceByUserID(ctx context.Context, userID int64) (*models.Balance, error) {
	query := `SELECT current_balance, total_withdrawn FROM balances WHERE user_id = $1`
	row := s.db.QueryRow(ctx, query, userID)
	s.logger.Debug("Getting balance", zap.Int64("user_id", userID))
	var balance models.Balance
	if err := row.Scan(&balance.CurrentBalance, &balance.TotalWithdrawn); err != nil {
		return nil, fmt.Errorf("error getting balance: %w", err)
	}
	s.logger.Debugf("Got balance: %v", balance)
	return &balance, nil
}

func (s *BalanceStorage) Withdraw(ctx context.Context, withdrawal *models.WithdrawalRequest) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
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

	updateBalanceQuery := `
        UPDATE balances
        SET current_balance = current_balance - $1, total_withdrawn = total_withdrawn + $1
        WHERE user_id = $2`
	_, err = tx.Exec(ctx, updateBalanceQuery, withdrawal.Amount, withdrawal.UserID)
	if err != nil {
		return fmt.Errorf("error updating balance: %w", err)
	}

	insertWithdrawalQuery := `
        INSERT INTO withdrawals (user_id, order_number, sum, processed_at)
        VALUES ($1, $2, $3, NOW())`
	_, err = tx.Exec(ctx, insertWithdrawalQuery, withdrawal.UserID, withdrawal.Order, withdrawal.Amount)
	if err != nil {
		return fmt.Errorf("error inserting withdrawal: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}

func (s *BalanceStorage) GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]models.Withdrawal, error) {
	query := `SELECT order_number, sum, processed_at FROM withdrawals WHERE user_id = $1` // ORDER BY processed_at
	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error db query of get withdrawals by user id: %w", err)
	}
	defer rows.Close()

	var withdrawals []models.Withdrawal
	for rows.Next() {
		var withdrawal models.Withdrawal
		if err := rows.Scan(&withdrawal.Order, &withdrawal.Amount, &withdrawal.ProcessedAt); err != nil {
			return nil, fmt.Errorf("error scan model withdrawal %w", err)
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error read rows by get withdrawal by user %w", err)
	}

	return withdrawals, nil
}

func (s *BalanceStorage) CreateBalanceTx(ctx context.Context, tx pgx.Tx, userID int64) error {
	query := "INSERT INTO balances (user_id, current_balance, total_withdrawn) VALUES ($1, $2, $3)"
	_, err := tx.Exec(ctx, query, userID, 0, 0)
	if err != nil {
		return fmt.Errorf("ошибка при создании баланса: %w", err)
	}
	return nil
}

func (s *BalanceStorage) CreateBalance(ctx context.Context, userID int64) error {
	query := "INSERT INTO balances (user_id, current_balance, total_withdrawn) VALUES ($1, $2, $3)"
	_, err := s.db.Exec(ctx, query, userID, 0, 0)
	if err != nil {
		return fmt.Errorf("ошибка при создании баланса: %w", err)
	}
	return nil
}

func (s *BalanceStorage) UpdateBalance(ctx context.Context, userID int64, amount float32) error {
	s.logger.Info("Updating balance", zap.Int64("user_id", userID))
	query := `
        UPDATE balances
        SET current_balance = current_balance + $1
        WHERE user_id = $2`
	_, err := s.db.Exec(ctx, query, amount, userID)
	if err != nil {
		return fmt.Errorf("error updating balance: %w", err)
	}
	s.logger.Info("Updated balance")
	return nil
}

func (s *BalanceStorage) UpdateBalanceTx(ctx context.Context, tx pgx.Tx, userID int64, accrual float32) error {
	query := "UPDATE balances SET current_balance = current_balance + $1 WHERE user_id = $2"
	_, err := tx.Exec(ctx, query, accrual, userID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении баланса: %w", err)
	}
	s.logger.Warnf("UpdateBalanceTx add to balance %v for user %v", accrual, userID)
	return nil
}
