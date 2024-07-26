package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/repositories/balancerepo"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/repositories/orderrepo"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/repositories/userrepo"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/config"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/handlers"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/services"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/storage"
	"github.com/Erlast/loyalty-gophermart.git/pkg/zaplog"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

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

func getAllTables(ctx context.Context, db *pgxpool.Pool) ([]string, error) {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
	`
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("could not scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", rows.Err())
	}

	return tables, nil
}

func main() {
	ctx := context.Background()

	logger := zaplog.InitLogger()
	logger.Infof("Logger started at %s", time.Now)
	defer func(Logger *zap.SugaredLogger) {
		err := Logger.Sync()
		if err != nil {
			log.Printf("Error syncing zap logger: %s", err)
		}
	}(logger)

	cfg := config.LoadConfig(logger)
	logger.Infof("Config: %v", cfg)

	fmt.Println("gophermart", cfg.DatabaseURI)

	db, err := storage.InitDB(ctx, cfg)
	if err != nil {
		fmt.Println("error initializing database", err)
		logger.Fatalf("Error initializing database: %s", err)
	}
	defer db.Close()
	logger.Infof("Database initialized")

	tables, err := getAllTables(ctx, db)
	if err != nil {
		fmt.Println("error getting all tables", err)
		log.Fatalf("Error retrieving tables: %v\n", err)
	}

	// Print the retrieved table names
	fmt.Println("Tables in the public schema:")
	for _, table := range tables {
		fmt.Println(table)
	}

	time.Sleep(2 * time.Second)

	// Инициализация сервисов
	userStorage := userrepo.NewUserStorage(db, logger)
	orderStorage := orderrepo.NewOrderStorage(db, logger)
	balanceStorage := balancerepo.NewBalanceStorage(db, logger)
	accrualService := services.NewAccrualService(cfg.AccrualSystemAddress, logger)
	userService := services.NewUserService(userStorage, balanceStorage, logger)
	orderService := services.NewOrderService(orderStorage, balanceStorage, accrualService, logger)
	balanceService := services.NewBalanceService(balanceStorage, logger)

	// Инициализация роутера
	router := chi.NewRouter()
	logger.Infof("Initializing router")

	// Регистрация маршрутов
	handlers.RegisterRoutes(ctx, router, userService, orderService, balanceService, logger)
	logger.Infof("Routes registered")

	// Запуск фоновой горутины для обновления статусов заказов
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				logger.Error("orderService.UpdateOrderStatuses called every 1 minute")
				if err := orderService.UpdateOrderStatuses(ctx); err != nil {
					logger.Error("Error updating order statuses", zap.Error(err))
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	logger.Infof("Update order statuses running")

	// Настройка и запуск сервера
	srv := &http.Server{
		Addr:    config.GetConfig().RunAddress,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("ListenAndServe failed", zap.Error(err))
		}
	}()
	logger.Info("Server is running", zap.String("address", config.GetConfig().RunAddress))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5 секунд на завершение
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
