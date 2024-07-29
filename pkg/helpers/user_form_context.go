package helpers

import (
	"context"
	"fmt"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/middleware"

	"go.uber.org/zap"
)

type FromContext interface {
	GetUserID(ctx context.Context, logger *zap.SugaredLogger) (int64, error)
}

type UserFormContext struct{}

func (m *UserFormContext) GetUserID(ctx context.Context, logger *zap.SugaredLogger) (int64, error) {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("error getting user id from context: %w", err)
	}
	logger.Info("Got user ID from context", zap.Int64("UserID", userID))
	return userID, nil
}
