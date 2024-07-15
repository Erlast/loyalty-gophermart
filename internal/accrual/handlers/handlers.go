package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
	"github.com/Erlast/loyalty-gophermart.git/pkg/zaplog"
)

func GetAccrualOrderHandler(
	_ context.Context,
	res http.ResponseWriter,
	req *http.Request,
	store storage.Storage,
) {
	orderNumber := chi.URLParam(req, "number")

	order, err := store.GetByOrderNumber(req.Context(), orderNumber)

	if err != nil {
		http.Error(res, "Not found", http.StatusNoContent)
		return
	}

	data, err := json.Marshal(order)
	if err != nil {
		zaplog.Logger.Errorf("failed to marshal result: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	_, err = res.Write(data)
	if err != nil {
		zaplog.Logger.Errorf("can't write body %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}
}

func PostAccrualOrderHandler(
	ctx context.Context,
	res http.ResponseWriter,
	req *http.Request,
	store storage.Storage,
) {
	var bodyReq models.OrderItem

	err := prepareBody(req, res, &bodyReq)

	if err != nil {
		return
	}

	err = store.SaveOrderItems(ctx, bodyReq)

	if err != nil {
		var conflictErr *helpers.ConflictError
		if errors.As(err, &conflictErr) {
			res.WriteHeader(http.StatusConflict)
			return
		}
		zaplog.Logger.Errorf("failed to save goods: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusAccepted)
}

func PostAccrualGoodsHandler(
	ctx context.Context,
	res http.ResponseWriter,
	req *http.Request,
	store storage.Storage,
) {
	var bodyReq models.Goods

	err := prepareBody(req, res, &bodyReq)

	if err != nil {
		return
	}

	err = store.SaveGoods(ctx, bodyReq)

	if err != nil {
		var conflictErr *helpers.ConflictError
		if errors.As(err, &conflictErr) {
			res.WriteHeader(http.StatusConflict)
			return
		}
		zaplog.Logger.Errorf("failed to save goods: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func prepareBody(req *http.Request, res http.ResponseWriter, bodyReq models.Model) error {
	if req.Body == http.NoBody {
		http.Error(res, "Empty body!", http.StatusBadRequest)
		return errors.New("empty body")
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		zaplog.Logger.Errorf("failed to read the request body: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return fmt.Errorf("read request body: %w", err)
	}

	err = json.Unmarshal(body, &bodyReq)

	if err != nil {
		zaplog.Logger.Errorf("failed to marshal result: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return fmt.Errorf("marashal result: %w", err)
	}

	if err := bodyReq.Validate(); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return fmt.Errorf("validation result: %w", err)
	}

	return nil
}
