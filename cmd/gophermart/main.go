package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/config"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/handlers"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/services"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/storage"
	"github.com/Erlast/loyalty-gophermart.git/pkg/zaplog"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

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

	db, err := storage.InitDB(ctx, cfg)
	if err != nil {
		logger.Fatalf("Error initializing database: %s", err)
	}
	defer db.Close()
	logger.Infof("Database initialized")

	// Инициализация сервисов
	userStorage := storage.NewUserStorage(db, logger)
	orderStorage := storage.NewOrderStorage(db, logger)
	balanceStorage := storage.NewBalanceStorage(db, logger)
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
		ticker := time.NewTicker(1 * time.Minute)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //nolint:mnd // 5 секунд на завершение
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
