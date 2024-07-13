package main

import (
	"context"
	"errors"
	"gofermart/internal/gofermart/config"
	"gofermart/internal/gofermart/handlers"
	"gofermart/internal/gofermart/middleware"
	"gofermart/internal/gofermart/services"
	"gofermart/internal/gofermart/storage"
	"gofermart/pkg/zaplog"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	zaplog.InitLogger()
	defer func(Logger *zap.SugaredLogger) {
		err := Logger.Sync()
		if err != nil {
			log.Printf("Error syncing zap logger: %s", err)
		}
	}(zaplog.Logger)

	cfg := config.LoadConfig()
	zaplog.Logger.Infof("Config: %v", cfg)

	err := storage.InitDB(ctx, cfg)
	if err != nil {
		zaplog.Logger.Fatalf("Error initializing database: %s", err)
	}
	defer storage.DB.Close()

	// Применение миграций
	err = storage.ApplyMigrations(ctx, "migrations/gofermart")
	if err != nil {
		zaplog.Logger.Fatal("Failed to apply migrations", zap.Error(err))
	}

	// Инициализация сервисов
	userStorage := storage.NewUserStorage(storage.DB)
	orderStorage := storage.NewOrderStorage(storage.DB)
	balanceStorage := storage.NewBalanceStorage(storage.DB)
	accrualService := services.NewAccrualService(cfg.AccrualSystemAddress)
	userService := services.NewUserService(userStorage)
	orderService := services.NewOrderService(orderStorage, accrualService)
	balanceService := services.NewBalanceService(balanceStorage)

	// Инициализация роутера
	router := chi.NewRouter()
	router.Use(middleware.AuthMiddleware(zaplog.Logger))

	// Регистрация маршрутов
	handlers.RegisterRoutes(ctx, router, userService, orderService, balanceService, zaplog.Logger)

	// Запуск фоновой горутины для обновления статусов заказов
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := orderService.UpdateOrderStatuses(ctx); err != nil {
					zaplog.Logger.Error("Error updating order statuses", zap.Error(err))
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Настройка и запуск сервера
	srv := &http.Server{
		Addr:    config.GetConfig().RunAddress,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zaplog.Logger.Fatal("ListenAndServe failed", zap.Error(err))
		}
	}()
	zaplog.Logger.Info("Server is running", zap.String("address", config.GetConfig().RunAddress))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	zaplog.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //nolint:mnd // 5 секунд на завершение
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zaplog.Logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	zaplog.Logger.Info("Server exiting")
}
