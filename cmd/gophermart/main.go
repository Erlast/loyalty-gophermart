package main

import (
	"github.com/go-chi/chi/v5"
	"gofermart/cmd/gophermart/internal/handlers"
	"gofermart/cmd/gophermart/pkg/config"
	"gofermart/cmd/gophermart/pkg/zaplog"
	"net/http"
)

func main() {
	zaplog.InitLogger()
	defer zaplog.Logger.Sync()

	cfg := config.LoadConfig()
	zaplog.Logger.Infof("Config: %v", cfg)

	// Create a new Chi router
	r := chi.NewRouter()

	// POST /api/user/register — регистрация пользователя
	r.Post("/api/user/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterUser()
	})

	// POST /api/user/login — аутентификация пользователя
	r.Post("/api/user/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginUser()
	})

	// POST /api/user/orders — загрузка пользователем номера заказа для расчёта
	r.Post("/api/user/orders", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoadOrders()
	})

	// GET /api/user/orders — получение списка загруженных пользователем номеров заказов и их статусов
	r.Get("/api/user/orders", func(w http.ResponseWriter, r *http.Request) {
		handlers.ListOrders()
	})

	// GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя
	r.Get("/api/user/balance", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetBalance()
	})

	// POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта
	r.Post("/api/user/balance/withdraw", func(w http.ResponseWriter, r *http.Request) {
		handlers.Withdraw()
	})

	// GET /api/user/withdrawals — получение информации о выводе средств с накопительного счёта
	r.Get("/api/user/withdrawals", func(w http.ResponseWriter, r *http.Request) {
		handlers.Withdrawals()
	})

	// Define your routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Gofermart!"))
	})

	// Start the HTTP server
	http.ListenAndServe(":8080", r)
}
