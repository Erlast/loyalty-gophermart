package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"
	models2 "github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func GetAccrualOrderHandler(_ context.Context, res http.ResponseWriter, req *http.Request, store *storage.AccrualStorage, log *zap.SugaredLogger) {
	orderNumber := chi.URLParam(req, "number")

	order, err := store.GetByOrderNumber(req.Context(), orderNumber)

	if err != nil {
		http.Error(res, "Not found", http.StatusNoContent)
		return
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Errorf("failed to marshal result: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	_, err = res.Write(data)
	if err != nil {
		log.Errorf("can't write body %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}
}

func PostAccrualOrderHandler(_ context.Context, res http.ResponseWriter, req *http.Request, store *storage.AccrualStorage, log *zap.SugaredLogger) {
	if req.Body == http.NoBody {
		http.Error(res, "Empty body!", http.StatusBadRequest)
		return
	}
	var bodyReq models2.OrderItem

	body, err := io.ReadAll(req.Body)

	if err != nil {
		log.Errorf("failed to read the request body: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &bodyReq)

	if err != nil {
		log.Errorf("failed to marshal result: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	if err := bodyReq.Validate(); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = store.SaveOrderItems(req.Context(), bodyReq)

	if err != nil {
		var conflictErr *helpers.ConflictError
		if errors.As(err, &conflictErr) {
			res.WriteHeader(http.StatusConflict)
			return
		}
		log.Errorf("failed to save order: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusAccepted)

}

func PostAccrualGoodsHandler(_ context.Context, res http.ResponseWriter, req *http.Request, store *storage.AccrualStorage, log *zap.SugaredLogger) {
	if req.Body == http.NoBody {
		http.Error(res, "Empty body!", http.StatusBadRequest)
		return
	}
	var bodyReq models2.Goods

	body, err := io.ReadAll(req.Body)

	if err != nil {
		log.Errorf("failed to read the request body: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &bodyReq)

	if err != nil {
		log.Errorf("failed to marshal result: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	if err := bodyReq.Validate(); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = store.SaveGoods(req.Context(), bodyReq)

	if err != nil {
		var conflictErr *helpers.ConflictError
		if errors.As(err, &conflictErr) {
			res.WriteHeader(http.StatusConflict)
			return
		}
		log.Errorf("failed to save order: %v", err)
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)

}
