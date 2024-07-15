package routes

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/handlers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
)

func NewAccrualRouter(
	ctx context.Context,
	store *storage.AccrualStorage,
	log *zap.SugaredLogger,
) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/orders/{number}", func(res http.ResponseWriter, req *http.Request) {
		handlers.GetAccrualOrderHandler(ctx, res, req, store, log)
	})

	r.Post("/api/orders", func(res http.ResponseWriter, req *http.Request) {
		handlers.PostAccrualOrderHandler(ctx, res, req, store, log)
	})

	r.Post("/api/goods", func(res http.ResponseWriter, req *http.Request) {
		handlers.PostAccrualGoodsHandler(ctx, res, req, store, log)
	})

	return r
}