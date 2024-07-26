package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"

	"github.com/Erlast/loyalty-gophermart.git/pkg/zaplog"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/components"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/config"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/routes"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
)

func main() {
	ctx := context.Background()
	newLogger := zaplog.InitLogger()
	cfg := config.ParseFlags(newLogger)

	store, err := storage.NewAccrualStorage(ctx, cfg)
	if err != nil {
		newLogger.Fatalf("Unable to create storage %v: ", err)
	}

	fmt.Println("accrual", cfg.DatabaseURI)

	go components.OrderProcessing(ctx, store, newLogger)

	r := routes.NewAccrualRouter(ctx, store, newLogger)

	newLogger.Infof("Start running server. Address: %s, db: %s", cfg.RunAddress, cfg.DatabaseURI)

	srv := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			newLogger.Fatal("Running server fail", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	newLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		newLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	newLogger.Info("Server exiting")
}
