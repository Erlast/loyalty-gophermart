package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Erlast/loyalty-gophermart.git/pkg/opensearch"
	"io"
	"net/http"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
	"github.com/go-chi/chi/v5"
)

func GetAccrualOrderHandler(
	_ context.Context,
	res http.ResponseWriter,
	req *http.Request,
	store storage.Storage,
	logger *opensearch.Logger,
) {
	orderNumber := chi.URLParam(req, "number")

	order, err := store.GetByOrderNumber(req.Context(), orderNumber)

	if err != nil {
		http.Error(res, "Not found", http.StatusNoContent)
		return
	}

	data, err := json.Marshal(order)
	if err != nil {
		logger.SendLog("error", fmt.Sprintf("failed to marshal result: %v", err))
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	_, err = res.Write(data)
	if err != nil {
		logger.SendLog("error", fmt.Sprintf("can't write body %v", err))
		http.Error(res, "", http.StatusInternalServerError)
		return
	}
}

func PostAccrualOrderHandler(
	ctx context.Context,
	res http.ResponseWriter,
	req *http.Request,
	store storage.Storage,
	logger *opensearch.Logger,
) {
	var bodyReq models.OrderItem

	err := prepareBody(req, res, &bodyReq, logger)

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

		logger.SendLog("error", fmt.Sprintf("failed to save goods: %v", err))
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
	logger *opensearch.Logger,
) {
	var bodyReq models.Goods

	err := prepareBody(req, res, &bodyReq, logger)

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
		logger.SendLog("error", fmt.Sprintf("failed to save goods: %v", err))
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func prepareBody(req *http.Request, res http.ResponseWriter, bodyReq models.Model, logger *opensearch.Logger) error {
	if req.Body == http.NoBody {
		http.Error(res, "Empty body!", http.StatusBadRequest)
		return errors.New("empty body")
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		logger.SendLog("error", fmt.Sprintf("failed to read the request body: %v", err))
		http.Error(res, "", http.StatusInternalServerError)
		return fmt.Errorf("read request body: %w", err)
	}

	err = json.Unmarshal(body, &bodyReq)

	if err != nil {
		logger.SendLog("error", fmt.Sprintf("failed to marshal result: %v", err))
		http.Error(res, "", http.StatusInternalServerError)
		return fmt.Errorf("marashal result: %w", err)
	}

	if err := bodyReq.Validate(); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return fmt.Errorf("validation result: %w", err)
	}

	return nil
}
