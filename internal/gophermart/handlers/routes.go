package handlers

import (
	"context"
	"github.com/Erlast/loyalty-gophermart.git/pkg/helpers"
	"net/http"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/middleware"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/services"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func RegisterRoutes(
	ctx context.Context,
	r *chi.Mux,
	userService *services.UserService,
	orderService *services.OrderService,
	balanceService *services.BalanceService,
	logger *zap.SugaredLogger,
) {
	userHandler := NewUserHandler(userService, logger)
	orderHandler := NewOrderHandler(orderService, logger)
	balanceHandler := NewBalanceHandler(balanceService, logger)

	fromContext := new(helpers.UserFormContext)

	// POST /api/user/register — регистрация пользователя
	r.Post("/api/user/register", userHandler.Register)
	// POST /api/user/login — аутентификация пользователя
	r.Post("/api/user/login", userHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(logger))

		// загрузка пользователем номера заказа для расчёта балов
		r.Post("/api/user/orders", func(w http.ResponseWriter, r *http.Request) {
			orderHandler.LoadOrder(ctx, w, r)
		})

		// получение списка загруженных пользователем номеров заказов и их статусов
		r.Get("/api/user/orders", func(w http.ResponseWriter, r *http.Request) {
			orderHandler.ListOrders(ctx, w, r)
		})

		// получение текущего баланса счёта баллов лояльности пользователя
		r.Get("/api/user/balance", func(w http.ResponseWriter, r *http.Request) {
			balanceHandler.GetBalance(ctx, w, r, fromContext)
		})

		// запрос на списание баллов с накопительного счёта
		r.Post("/api/user/balance/withdraw", func(w http.ResponseWriter, r *http.Request) {
			balanceHandler.Withdraw(ctx, w, r)
		})

		// получение информации о выводе средств с накопительного счёта
		r.Get("/api/user/withdrawals", func(w http.ResponseWriter, r *http.Request) {
			balanceHandler.Withdrawals(ctx, w, r)
		})
	})
}
