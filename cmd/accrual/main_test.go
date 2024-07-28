package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/config"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/models"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/routes"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
)

func TestParseFlags(t *testing.T) {
	t.Setenv("RUN_ADDRESS", "localhost:9090")
	t.Setenv("DATABASE_URI", "postgres://user:password@localhost:5432/dbname")

	defer func() {
		err := os.Unsetenv("RUN_ADDRESS")
		if err != nil {
			t.Fatalf("failed to unset env var: %v", err)
		}
	}()
	defer func() {
		err := os.Unsetenv("DATABASE_URI")
		if err != nil {
			t.Fatalf("failed to unset env var: %v", err)
		}
	}()

	logger := zap.NewNop().Sugar()
	cfg := config.ParseFlags(logger)

	assert.Equal(t, "localhost:9090", cfg.RunAddress)
	assert.Equal(t, "postgres://user:password@localhost:5432/dbname", cfg.DatabaseURI)
}

func TestNewAccrualRouter(t *testing.T) {
	logger := zap.NewNop().Sugar()
	store := &storage.MockStorage{}

	router := routes.NewAccrualRouter(store, logger)

	ts := httptest.NewServer(router)
	defer ts.Close()
	store.On("GetByOrderNumber", mock.Anything, "1d21").Return(&models.Order{
		UUID:    "uuid",
		Status:  helpers.StatusRegistered,
		Accrual: 100.00,
	}, nil)
	resp, err := http.Get(ts.URL + "/api/orders/1d21")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	err = resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to close response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}
}
