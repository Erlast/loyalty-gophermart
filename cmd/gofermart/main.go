package main

import (
	"gofermart/internal/gofermart/config"
	"gofermart/internal/gofermart/handlers"
	"gofermart/internal/gofermart/middleware"
	"gofermart/internal/gofermart/services"
	"gofermart/internal/gofermart/storage"
	"gofermart/pkg/zaplog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	zaplog.InitLogger()
	defer zaplog.Logger.Sync()

	cfg := config.LoadConfig()
	zaplog.Logger.Infof("Config: %v", cfg)

	storage.InitDB(cfg)
	defer storage.DB.Close()

	// Применение миграций
	err := storage.ApplyMigrations("migrations/gofermart")
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

		// POST /api/user/orders — загрузка пользователем номера заказа для расчёта
		r.Post("/api/user/orders", orderHandler.LoadOrders)

		// GET /api/user/orders — получение списка загруженных пользователем номеров заказов и их статусов
		r.Get("/api/user/orders", orderHandler.ListOrders)

		// GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя
		r.Get("/api/user/balance", balanceHandler.GetBalance)

		// POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта
		r.Post("/api/user/balance/withdraw", balanceHandler.Withdraw)

		// GET /api/user/withdrawals — получение информации о выводе средств с накопительного счёта
		r.Get("/api/user/withdrawals", balanceHandler.Withdrawals)
	})

	// Define your routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Gofermart!"))
	})

	// Start the HTTP server
	http.ListenAndServe(":8080", r)
}
