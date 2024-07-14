package storage

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"loyalty-gophermart/internal/gofermart/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(ctx context.Context, cfg config.Config, migrationsDir embed.FS) error {
	if err := runMigrations(cfg.DatabaseURI, migrationsDir); err != nil {
		return fmt.Errorf("failed to run DB migrations: %w", err)
	}

	parsedConfig, err := pgxpool.ParseConfig(cfg.DatabaseURI)
	if err != nil {
		return fmt.Errorf("unable to parse cfg.DatabaseURI: %w", err)
	}

	DB, err = pgxpool.ConnectConfig(ctx, parsedConfig)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
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
