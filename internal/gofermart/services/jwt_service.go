package services

import (
	"fmt"
	"gofermart/internal/gofermart/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTClaim struct {
	jwt.StandardClaims
	UserID int64 `json:"user_id"`
}

func GenerateJWT(userID int64) (string, error) {
	claims := &JWTClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), //nolint:mnd // 72 часа
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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
