package middleware

import (
	"context"
	"errors"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/services/jwt"

	"net/http"
	"strings"

	"go.uber.org/zap"
)

type contextKey string

var contextKeyUserID = contextKey("userID")

func SetUserIDToContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, contextKeyUserID, userID)
}

func GetUserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(contextKeyUserID).(int64)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}

func AuthMiddleware(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			logger.Infof("Auth header: %s", authHeader)
			if authHeader == "" {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := jwt.ParseJWT(tokenStr)
			logger.Infof("Claims: %v", claims)
			if err != nil {
				logger.Error("Invalid token", zap.Error(err))
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			ctx := SetUserIDToContext(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
