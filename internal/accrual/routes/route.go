package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/handlers"
)

func NewAccrualRouter(
	store handlers.Storage,
	logger *zap.SugaredLogger,
) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/orders/{number}", func(res http.ResponseWriter, req *http.Request) {
		handlers.GetAccrualOrderHandler(res, req, store, logger)
	})

	r.Post("/api/orders", func(res http.ResponseWriter, req *http.Request) {
		handlers.PostAccrualOrderHandler(res, req, store, logger)
	})

	r.Post("/api/goods", func(res http.ResponseWriter, req *http.Request) {
		handlers.PostAccrualGoodsHandler(res, req, store, logger)
	})

	return r
}
