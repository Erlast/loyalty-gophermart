package helpers

import (
	"errors"
	"net/http"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/config"

	"go.uber.org/zap"
)

// GetUserIDFromContext извлекает userID из контекста и возвращает его, или ошибку, если не удалось извлечь.
func GetUserIDFromContext(r *http.Request, logger *zap.SugaredLogger) (int64, error) {
	userIDValue := r.Context().Value(config.UserIDContextKey)
	userID, ok := userIDValue.(int64)
	if !ok {
		logger.Error("Error getting userID from context", zap.Any("userIDValue", userIDValue))
		return 0, errors.New("error getting userID from context")
	}
	return userID, nil
}
