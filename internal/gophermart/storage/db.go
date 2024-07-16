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

var DB *pgxpool.Pool

func InitDB(ctx context.Context, cfg config.Config) error {
	// Запускаем миграции
	if err := runMigrations(cfg.DatabaseURI, migrationsDir); err != nil {
		return fmt.Errorf("не удалось выполнить миграции в базу данных: %w", err)
	}

	// Парсим URI базы данных
	parsedConfig, err := pgxpool.ParseConfig(cfg.DatabaseURI)
	if err != nil {
		return fmt.Errorf("не удалось разобрать URI базы данных: %w", err)
	}

	// Подключаемся к базе данных
	DB, err = pgxpool.NewWithConfig(ctx, parsedConfig)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	return nil
}

func runMigrations(dsn string, migrationsDir embed.FS) error {
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
