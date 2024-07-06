package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"gofermart/internal/gofermart/config"
	"os"
	"path/filepath"
)

var DB *pgxpool.Pool

func InitDB(cfg config.Config) error {
	parsedConfig, err := pgxpool.ParseConfig(cfg.DatabaseURI)
	if err != nil {
		return fmt.Errorf("Unable to parse cfg.DatabaseURI: %v\n", err)
	}

	DB, err = pgxpool.ConnectConfig(context.Background(), parsedConfig)
	if err != nil {
		return fmt.Errorf("Unable to connect to database: %v\n", err)
	}

	return nil
}

func CloseDB() {
	DB.Close()
}

func ApplyMigrations(migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(migrationsDir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", filePath, err)
		}

		_, err = DB.Exec(context.Background(), string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration file %s: %v", filePath, err)
		}
	}

	return nil
}
