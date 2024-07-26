package storage

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsDir embed.FS

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

func InitDB(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	// Запускаем миграции
	if err := runMigrations(cfg.DatabaseURI); err != nil {
		return nil, fmt.Errorf("не удалось выполнить миграции в базу данных: %w", err)
	}

	// Парсим URI базы данных
	parsedConfig, err := pgxpool.ParseConfig(cfg.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("не удалось разобрать URI базы данных: %w", err)
	}

	// Подключаемся к базе данных
	db, err := pgxpool.NewWithConfig(ctx, parsedConfig)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	exist, err := tableExists(ctx, db, "users")
	if err != nil {
		return nil, fmt.Errorf("error checking if table exists: %w", err)
	}
	if !exist {
		return nil, fmt.Errorf("table users doesn't exists ")
	} else {
		return db, fmt.Errorf("tablle exists")
	}

	//return db, nil
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
