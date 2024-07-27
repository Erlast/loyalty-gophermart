package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/config"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/repositories/balancerepo"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/repositories/userrepo"
	"github.com/Erlast/loyalty-gophermart.git/internal/gophermart/services"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestUserHandler_Login_StatusOK(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t).Sugar()
	timeUpdatedAt := time.Now().Add(-24 * time.Hour)
	timeCreatedAt := time.Now().Add(-48 * time.Hour)

	mockUserStorage := new(userrepo.MockUserStore)
	mockBalanceStorage := new(balancerepo.MockBalanceStore)

	mockUserStorage.On("GetUserByLogin", ctx, "admin10").Return(&models.User{
		ID:        1,
		Login:     "admin10",
		Password:  "$2b$12$Mv.HOsLLFy8MQGkIlaxI.u3ruZ/4C4JyKJamyAjm23C7uuLhFLfs6", // admin10
		UpdatedAt: timeUpdatedAt,
		CreatedAt: timeCreatedAt,
	}, nil)

	userService := services.NewUserService(mockUserStorage, mockBalanceStorage, logger)
	userHandler := NewUserHandler(userService, logger)

	body := strings.NewReader(`{"login": "admin10","password": "admin10"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", body)
	// Установка заголовка Content-Type
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	userHandler.Login(rr, req)

	mockUserStorage.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEmpty(t, rr.Body.String())
	assert.NotEmpty(t, rr.Header().Get("Authorization"))
}

func TestUserHandler_Login_StatusBadRequest(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	mockUserStorage := new(userrepo.MockUserStore)
	mockBalanceStorage := new(balancerepo.MockBalanceStore)

	userService := services.NewUserService(mockUserStorage, mockBalanceStorage, logger)
	userHandler := NewUserHandler(userService, logger)

	body := strings.NewReader(`{"login": "admin10"`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", body)
	// Установка заголовка Content-Type
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	userHandler.Login(rr, req)

	fmt.Println("Body:", rr.Body.String())
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserHandler_Login_StatusUnauthorized(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t).Sugar()
	timeUpdatedAt := time.Now().Add(-24 * time.Hour)
	timeCreatedAt := time.Now().Add(-48 * time.Hour)

	mockUserStorage := new(userrepo.MockUserStore)
	mockBalanceStorage := new(balancerepo.MockBalanceStore)

	mockUserStorage.On("GetUserByLogin", ctx, "admin10").Return(&models.User{
		ID:        1,
		Login:     "admin10",
		Password:  "passwordNotАdmin10",
		UpdatedAt: timeUpdatedAt,
		CreatedAt: timeCreatedAt,
	}, nil)

	userService := services.NewUserService(mockUserStorage, mockBalanceStorage, logger)
	userHandler := NewUserHandler(userService, logger)

	body := strings.NewReader(`{"login": "admin10","password": "admin10"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", body)
	// Установка заголовка Content-Type
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	userHandler.Login(rr, req)
	fmt.Println("Secret", config.GetConfig().JWTSecret)
	mockUserStorage.AssertExpectations(t)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
