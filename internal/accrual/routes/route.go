package routes

import (
	"context"
	"net/http"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/handlers"
)

type Storage interface {
	GetByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error)
	SaveOrderItems(ctx context.Context, items models.OrderItem) error
	SaveGoods(ctx context.Context, goods models.Goods) error
	GetRegisteredOrders(ctx context.Context) ([]int64, error)
	FetchRewardRules(ctx context.Context) ([]models.Goods, error)
	UpdateOrderStatus(ctx context.Context, orderNumber int64, status string) error
	FetchProducts(ctx context.Context, orderID int64) ([]models.Items, error)
	SaveOrderPoints(ctx context.Context, orderID int64, points []float32) error
}

func NewAccrualRouter(
	store Storage,
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
