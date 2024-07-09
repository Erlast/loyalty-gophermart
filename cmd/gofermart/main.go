package main

import (
	"context"
	"gofermart/internal/gofermart/config"
	"gofermart/internal/gofermart/handlers"
	"gofermart/internal/gofermart/middleware"
	"gofermart/internal/gofermart/services"
	"gofermart/internal/gofermart/storage"
	"gofermart/pkg/zaplog"
	"log"
	"net/http"

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

	// Create a new Chi router
	r := chi.NewRouter()

	// Инициализация сервисов
	userStorage := storage.NewUserStorage(storage.DB)
	orderStorage := storage.NewOrderStorage(storage.DB)
	balanceStorage := storage.NewBalanceStorage(storage.DB)
	userService := services.NewUserService(userStorage)
	orderService := services.NewOrderService(orderStorage)
	balanceService := services.NewBalanceService(balanceStorage)

	// инициализация обработчиков
	userHandler := handlers.NewUserHandler(userService, zaplog.Logger)
	orderHandler := handlers.NewOrderHandler(orderService, zaplog.Logger)
	balanceHandler := handlers.NewBalanceHandler(balanceService, zaplog.Logger)

	// POST /api/user/register — регистрация пользователя
	r.Post("/api/user/register", userHandler.Register)
	// POST /api/user/login — аутентификация пользователя
	r.Post("/api/user/login", userHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(zaplog.Logger))

		// загрузка пользователем номера заказа для расчёта балов
		r.Post("/api/user/orders", func(w http.ResponseWriter, r *http.Request) {
			orderHandler.LoadOrder(ctx, w, r)
		})

		// получение списка загруженных пользователем номеров заказов и их статусов
		r.Get("/api/user/orders", func(w http.ResponseWriter, r *http.Request) {
			orderHandler.ListOrders(ctx, w, r)
		})

		// запрос на списание баллов с накопительного счёта
		r.Post("/api/user/balance/withdraw", func(w http.ResponseWriter, r *http.Request) {
			balanceHandler.Withdraw(ctx, w, r)
		})

		// получение информации о выводе средств с накопительного счёта
		r.Get("/api/user/withdrawals", func(w http.ResponseWriter, r *http.Request) {
			balanceHandler.Withdrawals(ctx, w, r)
		})

		// получение текущего баланса счёта баллов лояльности пользователя
		r.Get("/api/user/balance", balanceHandler.GetBalance)
	})

	// Define your routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Welcome to Gofermart!"))
		if err != nil {
			zaplog.Logger.Errorf("Error writing response: %s", err)
		}
	})

	// Start the HTTP server
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		zaplog.Logger.Fatal("Error starting server", zap.Error(err))
	}
}
