package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
)

var limiter = rate.NewLimiter(1, 5)

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

func GetAccrualOrderHandler(
	res http.ResponseWriter,
	req *http.Request,
	store Storage,
	logger *zap.SugaredLogger,
) {
	if !limiter.Allow() {
		http.Error(res, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	orderNumber := chi.URLParam(req, "number")

	order, err := store.GetByOrderNumber(req.Context(), orderNumber)

	if err != nil {
		http.Error(res, "Not found", http.StatusNoContent)
		return
	}

	data, err := json.Marshal(order)
	if err != nil {
		logger.Errorf("failed to marshal result: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	_, err = res.Write(data)
	if err != nil {
		logger.Errorf("can't write body %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}
}

func PostAccrualOrderHandler(
	res http.ResponseWriter,
	req *http.Request,
	store Storage,
	logger *zap.SugaredLogger,
) {
	var bodyReq models.OrderItem

	err := prepareBody(req, res, &bodyReq, logger)

	if err != nil {
		return
	}

	err = store.SaveOrderItems(req.Context(), bodyReq)

	if err != nil {
		var conflictErr *helpers.ConflictError
		if errors.As(err, &conflictErr) {
			res.WriteHeader(http.StatusConflict)
			return
		}
		logger.Errorf("failed to save goods: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusAccepted)
}

func PostAccrualGoodsHandler(
	res http.ResponseWriter,
	req *http.Request,
	store Storage,
	logger *zap.SugaredLogger,
) {
	var bodyReq models.Goods

	err := prepareBody(req, res, &bodyReq, logger)

	if err != nil {
		return
	}

	err = store.SaveGoods(req.Context(), bodyReq)

	if err != nil {
		var conflictErr *helpers.ConflictError
		if errors.As(err, &conflictErr) {
			res.WriteHeader(http.StatusConflict)
			return
		}
		logger.Errorf("failed to save goods: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func prepareBody(req *http.Request, res http.ResponseWriter, bodyReq models.Model, logger *zap.SugaredLogger) error {
	if req.Body == http.NoBody {
		http.Error(res, "Empty body!", http.StatusBadRequest)
		return errors.New("empty body")
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		logger.Errorf("failed to read the request body: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return fmt.Errorf("read request body: %w", err)
	}

	err = json.Unmarshal(body, &bodyReq)

	if err != nil {
		logger.Errorf("failed to marshal result: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return fmt.Errorf("marashal result: %w", err)
	}

	if err := bodyReq.Validate(); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return fmt.Errorf("validation result: %w", err)
	}

	return nil
}
