package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/services"
	"github.com/Erlast/loyalty-gophermart.git/pkg/helpers"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type BalanceHandler struct {
	service *services.BalanceService
	logger  *zap.SugaredLogger
}

const ErrorGettingUserIDFromContext = "error_getting_user_id_from_context: %v"

func NewBalanceHandler(service *services.BalanceService, logger *zap.SugaredLogger) *BalanceHandler {
	return &BalanceHandler{service: service, logger: logger}
}

func (h *BalanceHandler) GetBalance(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserIDFromContext(r, h.logger)
	if err != nil {
		h.logger.Errorf(ErrorGettingUserIDFromContext, err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	balance, err := h.service.GetBalanceByUserID(ctx, userID)
	if err != nil {
		h.logger.Error("Error getting balance", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, balance)
}

func (h *BalanceHandler) Withdraw(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserIDFromContext(r, h.logger)
	if err != nil {
		h.logger.Error("Unauthorized user", zap.Error(err))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var withdrawal models.WithdrawalRequest
	if err := json.NewDecoder(r.Body).Decode(&withdrawal); err != nil {
		h.logger.Error("Error decoding request body", zap.Error(err))
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	withdrawal.UserID = userID

	err = h.service.Withdraw(ctx, &withdrawal)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInsufficientBalance):
			http.Error(w, "Insufficient balance", http.StatusPaymentRequired)
		case errors.Is(err, services.ErrInvalidOrderNumber):
			http.Error(w, InvalidOrderFormatMsg, http.StatusUnprocessableEntity)
		default:
			h.logger.Error("Error withdrawing balance", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BalanceHandler) Withdrawals(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserIDFromContext(r, h.logger)
	if err != nil {
		h.logger.Error(ErrorGettingUserIDFromContext, zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	withdrawals, err := h.service.GetWithdrawalsByUserID(ctx, userID)
	if err != nil {
		h.logger.Error("Error getting withdrawals", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	render.JSON(w, r, withdrawals)
}
