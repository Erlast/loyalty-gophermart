package helpers

import (
	"gofermart/internal/gofermart/config"
	"net/http"

	"go.uber.org/zap"
)

// GetUserIDFromContext извлекает userID из контекста и возвращает его, или ошибку, если не удалось извлечь.
func GetUserIDFromContext(r *http.Request, logger *zap.SugaredLogger) (int64, error) {
	userIDValue := r.Context().Value(config.UserIDContextKey)
	userID, ok := userIDValue.(int64)
	if !ok {
		logger.Error("Error getting userID from context", zap.Any("userIDValue", userIDValue))
		return 0, http.ErrNoLocation
	}
	return userID, nil
}