package balance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/middleware"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/services/balance"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/Erlast/loyalty-gophermart.git/pkg/helpers"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type BalanceHandler struct {
	service *balance.BalanceService
	logger  *zap.SugaredLogger
}

const ErrorGettingUserIDFromContext = "error_getting_user_id_from_context: %v"

func NewBalanceHandler(service *balance.BalanceService, logger *zap.SugaredLogger) *BalanceHandler {
	return &BalanceHandler{service: service, logger: logger}
}

func (h *BalanceHandler) GetBalance(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	fromContext helpers.FromContext,
) {
	userID, err := fromContext.GetUserID(r.Context(), h.logger)
	if err != nil {
		h.logger.Errorf(ErrorGettingUserIDFromContext, err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	balanceByUser, err := h.service.GetBalanceByUserID(ctx, userID)
	if err != nil {
		h.logger.Error("Error getting balance", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	h.logger.Debugf("Get Balance from handler: %v", balanceByUser)

	render.JSON(w, r, balanceByUser)
}

func (h *BalanceHandler) Withdraw(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Error("Unauthorized user", zap.Error(err))
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	var withdrawal models.WithdrawalRequest
	if err := json.NewDecoder(r.Body).Decode(&withdrawal); err != nil {
		h.logger.Error("Error decoding request body", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	withdrawal.UserID = userID

	err = h.service.Withdraw(ctx, &withdrawal)
	if err != nil {
		switch {
		case errors.Is(err, balance.ErrInsufficientBalance):
			http.Error(w, "", http.StatusPaymentRequired)
		case errors.Is(err, balance.ErrInvalidOrderNumber):
			http.Error(w, "", http.StatusUnprocessableEntity)
		default:
			fmt.Println("Error withdrawing balance:", err)
			h.logger.Error("Error withdrawing balance", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BalanceHandler) Withdrawals(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	fromContext helpers.FromContext,
) {
	userID, err := fromContext.GetUserID(r.Context(), h.logger)
	if err != nil {
		h.logger.Error(ErrorGettingUserIDFromContext, zap.Error(err))
		fmt.Println("error getting userID from context", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	withdrawals, err := h.service.GetWithdrawalsByUserID(ctx, userID)
	if err != nil {
		h.logger.Error("Error getting withdrawals", zap.Error(err))
		fmt.Println("error getting withdrawals", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	render.JSON(w, r, withdrawals)
}
