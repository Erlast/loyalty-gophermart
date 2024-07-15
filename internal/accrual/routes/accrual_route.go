package routes

import (
	"context"
	"net/http"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/handlers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
	"github.com/go-chi/chi/v5"
)

func NewAccrualRouter(
	ctx context.Context,
	store storage.Storage,
) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/orders/{number}", func(res http.ResponseWriter, req *http.Request) {
		handlers.GetAccrualOrderHandler(ctx, res, req, store)
	})

	r.Post("/api/orders", func(res http.ResponseWriter, req *http.Request) {
		handlers.PostAccrualOrderHandler(ctx, res, req, store)
	})

	r.Post("/api/goods", func(res http.ResponseWriter, req *http.Request) {
		handlers.PostAccrualGoodsHandler(ctx, res, req, store)
	})

	return r
}
