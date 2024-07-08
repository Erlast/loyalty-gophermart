package routes

import (
	"context"
	"github.com/Erlast/loyalty-gophermart.git/internal/handlers"
	"github.com/Erlast/loyalty-gophermart.git/internal/storage"
	"go.uber.org/zap"

	"github.com/Erlast/loyalty-gophermart.git/internal/config"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewAccrualRouter(ctx context.Context, _ *config.Cfg, store *storage.AccrualStorage, log *zap.SugaredLogger) *chi.Mux {
	r := chi.NewRouter()

	//r.Use(func(h http.Handler) http.Handler {
	//	return middlewares.AuthMiddleware(h, logger, conf)
	//})
	//r.Use(func(h http.Handler) http.Handler {
	//	return middlewares.WithLogging(h, logger)
	//})
	//r.Use(func(h http.Handler) http.Handler {
	//	return middlewares.GzipMiddleware(h, logger)
	//})

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
