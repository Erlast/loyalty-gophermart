package storage

import (
	"context"
	"fmt"
	"gofermart/internal/gofermart/config"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(cfg config.Config, ctx context.Context) error {
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

func ApplyMigrations(migrationsDir string, ctx context.Context) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(migrationsDir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filePath, err)
		}

		_, err = DB.Exec(ctx, string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration file %s: %w", filePath, err)
		}
	}

	return nil
}
