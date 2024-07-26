package helpers

import (
	"errors"
	"net/http"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/config"

	"go.uber.org/zap"
)

type FromContext interface {
	GetUserID(r *http.Request, logger *zap.SugaredLogger) (int64, error)
}

type UserFormContext struct{}

// GetUserID GetUserIDFromContext извлекает userID из контекста и возвращает его, или ошибку, если не удалось извлечь.
func (m *UserFormContext) GetUserID(r *http.Request, logger *zap.SugaredLogger) (int64, error) {
	userIDValue := r.Context().Value(config.UserIDContextKey)
	logger.Infof("Getting UserID from context %v", userIDValue)
	userID, ok := userIDValue.(int64)
	if !ok {
		logger.Error("Error getting userID from context", zap.Any("userIDValue", userIDValue))
		return 0, errors.New("error getting userID from context")
	}
	logger.Info("Got user ID from context", zap.Int64("UserID", userID))
	return userID, nil
}

// GetUserIDFromContext извлекает userID из контекста и возвращает его, или ошибку, если не удалось извлечь.
func GetUserIDFromContext(r *http.Request, logger *zap.SugaredLogger) (int64, error) {
	userIDValue := r.Context().Value(config.UserIDContextKey)
	logger.Infof("Getting UserID from context %v", userIDValue)
	userID, ok := userIDValue.(int64)
	if !ok {
		logger.Error("Error getting userID from context", zap.Any("userIDValue", userIDValue))
		return 0, errors.New("error getting userID from context")
	}
	logger.Info("Got user ID from context", zap.Int64("UserID", userID))
	return userID, nil
}
