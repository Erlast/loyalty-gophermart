package handlers

import (
	"encoding/json"
	"gofermart/internal/gofermart/models"
	"gofermart/internal/gofermart/services"
	"gofermart/pkg/helpers"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type BalanceHandler struct {
	service *services.BalanceService
	logger  *zap.SugaredLogger
}

func NewBalanceHandler(service *services.BalanceService, logger *zap.SugaredLogger) *BalanceHandler {
	return &BalanceHandler{service: service, logger: logger}
}

func (h *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserIDFromContext(r, h.logger)
	if err != nil {
		h.logger.Errorf("Error getting user id from context: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	balance, err := h.service.GetBalanceByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error("Error getting balance", zap.Error(err))
		http.Error(w, "Error getting balance", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, balance)
}

func (h *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserIDFromContext(r, h.logger)
	if err != nil {
		h.logger.Errorf("Error getting user id from context: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var withdrawal models.WithdrawalRequest
	if err := json.NewDecoder(r.Body).Decode(&withdrawal); err != nil {
		h.logger.Error("Error decoding request body", zap.Error(err))
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	withdrawal.UserID = userID

	if err := h.service.Withdraw(r.Context(), &withdrawal); err != nil {
		h.logger.Error("Error withdrawing balance", zap.Error(err))
		http.Error(w, "Error withdrawing balance: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BalanceHandler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserIDFromContext(r, h.logger)
	if err != nil {
		h.logger.Errorf("Error getting user id from context: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	withdrawals, err := h.service.GetWithdrawalsByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error("Error getting withdrawals", zap.Error(err))
		http.Error(w, "Error getting withdrawals", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, withdrawals)
}
