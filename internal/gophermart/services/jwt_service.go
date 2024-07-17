package services

import (
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/config"
	"github.com/dgrijalva/jwt-go"
)

type JWTClaim struct {
	jwt.StandardClaims
	UserID int64 `json:"user_id"`
}

func GenerateJWT(userID int64, logger *zap.SugaredLogger) (string, error) {
	logger.Infof("User ID: %v", userID)
	claims := &JWTClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), //nolint:mnd // 72 часа
		},
	}

	logger.Infof("JWT: %v", claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	logger.Infof("token: %v", token)
	signedJWT, err := token.SignedString([]byte(config.GetConfig().JWTSecret))
	if err != nil {
		return "", fmt.Errorf("error while signing JWT: %w", err)
	}
	return signedJWT, nil
}

// ParseJWT jwt.ParseWithClaims — это метод из библиотеки github.com/dgrijalva/jwt-go
// func(token *jwt.Token) (interface{}, error) — это функция, которая предоставляется для валидации токена.
// Она должна возвращать секретный ключ, используемый для подписи токена.
func ParseJWT(tokenStr string) (*JWTClaim, error) {
	claims := &JWTClaim{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("error while parsing JWT: %w", err)
	}

	return claims, nil
}
