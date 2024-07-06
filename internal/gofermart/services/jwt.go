package services

import (
	"gofermart/internal/gofermart/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTClaim struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(userID int64) (string, error) {
	claims := &JWTClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetConfig().JWTSecret))
}

func ParseJWT(tokenStr string) (*JWTClaim, error) {
	claims := &JWTClaim{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
