package storage

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/config"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
)

type Storage interface {
	GetByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error)
	SaveOrderItems(ctx context.Context, items models.OrderItem) error
	SaveGoods(ctx context.Context, goods models.Goods) error
	GetRegisteredOrders(ctx context.Context) ([]int64, error)
	FetchRewardRules(ctx context.Context) ([]models.Goods, error)
	UpdateOrderStatus(ctx context.Context, orderNumber int64, status string) error
	FetchProducts(ctx context.Context, orderID int64) ([]models.Items, error)
	SaveOrderPoints(ctx context.Context, orderID int64, points []int64) error
}

//go:embed migrations/*.sql
var migrationsDir embed.FS

type AccrualStorage struct {
	db *pgxpool.Pool
}

func NewAccrualStorage(ctx context.Context, cfg *config.Cfg) (*AccrualStorage, error) {
	if cfg.DatabaseURI == "" {
		return nil, errors.New("database uri is empty")
	}

	if err := runMigrations(cfg.DatabaseURI); err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}
	conn, err := initPool(ctx, cfg.DatabaseURI)

	if err != nil {
		return nil, fmt.Errorf("unable to connect database: %w", err)
	}

	return &AccrualStorage{db: conn}, nil
}

func (store *AccrualStorage) GetByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
	var order models.Order
	err := store.db.QueryRow(ctx, "SELECT uuid, status,accrual FROM orders WHERE uuid = $1", orderNumber).Scan(
		&order.UUID,
		&order.Status,
		&order.Accrual,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("order not found %w", err)
		}
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	return &order, nil
}

func (store *AccrualStorage) SaveOrder(ctx context.Context, orderNumber string) (int64, error) {
	var orderID int64
	sqlString := "INSERT INTO orders(uuid, status, accrual, uploaded_at) VALUES ($1, $2, $3, $4) RETURNING id"
	err := store.db.QueryRow(ctx, sqlString, orderNumber, helpers.StatusRegistered, 0, time.Now()).Scan(&orderID)

	if err != nil {
		var pgsErr *pgconn.PgError
		if errors.As(err, &pgsErr) && pgsErr.Code == pgerrcode.UniqueViolation {
			return 0, &helpers.ConflictError{
				OrderNumber: orderNumber,
				Err:         err,
			}
		}
		return 0, fmt.Errorf("unable to save order: %w", err)
	}
	return orderID, nil
}

func (store *AccrualStorage) SaveOrderItems(ctx context.Context, order models.OrderItem) error {
	id, err := store.SaveOrder(ctx, order.UUID)
	if err != nil {
		return fmt.Errorf("unable to save order: %w", err)
	}

	batch := &pgx.Batch{}
	stmt := "INSERT INTO order_items(order_id, price, description) VALUES (@order_id,@price,@description)"

	for _, item := range order.Goods {
		args := pgx.NamedArgs{"order_id": id, "price": item.Price, "description": item.Description}
		batch.Queue(stmt, args)
	}

	tx, err := store.db.Begin(ctx)

	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if e := tx.Rollback(ctx); e != nil {
			err = fmt.Errorf("failed to rollback the transaction: %w", e)
			return
		}
	}()

	results := tx.SendBatch(ctx, batch)

	defer func() {
		if e := results.Close(); e != nil {
			err = fmt.Errorf("closing batch results error: %w", e)
			return
		}

		if e := tx.Commit(ctx); e != nil {
			err = fmt.Errorf("unable to commit: %w", e)
			return
		}
	}()

	return nil
}

func (store *AccrualStorage) SaveGoods(ctx context.Context, goods models.Goods) error {
	sqlString := "INSERT INTO accrual_rules(match, reward, reward_type) VALUES ($1, $2, $3)"
	_, err := store.db.Exec(ctx, sqlString, goods.Match, goods.Reward, goods.RewardType)

	if err != nil {
		var pgsErr *pgconn.PgError
		if errors.As(err, &pgsErr) && pgsErr.Code == pgerrcode.UniqueViolation {
			return &helpers.ConflictError{
				OrderNumber: "0",
				Err:         err,
			}
		}
		return fmt.Errorf("unable to save goods: %w", err)
	}
	return nil
}

func (store *AccrualStorage) GetRegisteredOrders(ctx context.Context) ([]int64, error) {
	query := `Select id FROM orders WHERE status=$1`
	rows, err := store.db.Query(ctx, query, helpers.StatusRegistered)
	if err != nil {
		return nil, fmt.Errorf("unable to get registered orders: %w", err)
	}

	var orders []int64

	for rows.Next() {
		var orderID int64
		err = rows.Scan(&orderID)
		if err != nil {
			return nil, fmt.Errorf("unable to get order id: %w", err)
		}
		orders = append(orders, orderID)
	}

	return orders, nil
}

func (store *AccrualStorage) FetchRewardRules(ctx context.Context) ([]models.Goods, error) {
	var rules []models.Goods
	rows, err := store.db.Query(ctx, "SELECT match, reward, reward_type FROM accrual_rules")

	if err != nil {
		return nil, fmt.Errorf("can't get rules. %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Goods
		if err := rows.Scan(&r.Match, &r.Reward, &r.RewardType); err != nil {
			return nil, fmt.Errorf("can't parse rule. %w", err)
		}
		rules = append(rules, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't parse error row. %w", err)
	}
	return rules, nil
}

func (store *AccrualStorage) FetchProducts(ctx context.Context, orderID int64) ([]models.Items, error) {
	rows, err := store.db.Query(ctx, "SELECT description, price FROM order_items WHERE order_id = $1", orderID)
	if err != nil {
		return nil, fmt.Errorf("can't get products. %w", err)
	}
	defer rows.Close()

	var products []models.Items
	for rows.Next() {
		var p models.Items
		if err := rows.Scan(&p.Description, &p.Price); err != nil {
			return nil, fmt.Errorf("can't parse product. %w", err)
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't parse error row. %w", err)
	}
	return products, nil
}

func (store *AccrualStorage) SaveOrderPoints(ctx context.Context, orderID int64, points []int64) error {
	var totalPoints int64
	for _, p := range points {
		totalPoints += p
	}

	_, err := store.db.Exec(
		ctx,
		"UPDATE orders SET status=$1,accrual=$2 where id=$3",
		helpers.StatusProcessed,
		totalPoints,
		orderID,
	)
	if err != nil {
		return fmt.Errorf("can't update orders. %w", err)
	}

	return nil
}

func (store *AccrualStorage) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	_, err := store.db.Exec(ctx, "Update orders set status=$1 where id=$2", status, orderID)

	if err != nil {
		return fmt.Errorf("can't update order. %w", err)
	}
	return nil
}

func initPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the DSN: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping the DB: %w", err)
	}
	return pool, nil
}

func runMigrations(dsn string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}
	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
	}
	return nil
}