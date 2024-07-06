package middleware

import (
	"context"
	"go.uber.org/zap"
	"gofermart/pkg/jwt"
	"net/http"
	"strings"
)

func AuthMiddleware(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwt.ParseJWT(tokenStr)
			if err != nil {
				logger.Error("Invalid token", zap.Error(err))
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
