package middleware

import (
	"context"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/services/jwt"

	"net/http"
	"strings"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/config"

	"go.uber.org/zap"
)

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

			ctx := context.WithValue(r.Context(), config.UserIDContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
