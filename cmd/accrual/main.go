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

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/components"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/config"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/routes"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
	"github.com/Erlast/loyalty-gophermart.git/pkg/opensearch"
)

func main() {
	ctx := context.Background()

	newLogger, err := opensearch.NewOpenSearchLogger()

	if err != nil {
		fmt.Printf("Error creating logger: %s\n", err)
		return
	}
	defer func(Logger *zap.Logger) {
		err := Logger.Sync()
		if err != nil {
			fmt.Printf("Error closing logger: %s\n", err)
			return
		}
	}(newLogger.Logger)

	cfg := config.ParseFlags(newLogger)

	store, err := storage.NewAccrualStorage(ctx, cfg)
	if err != nil {
		newLogger.SendLog("fatal", fmt.Sprintf("Unable to create storage %v: ", err))
	}

	go components.OrderProcessing(ctx, store, newLogger)

	r := routes.NewAccrualRouter(ctx, store, newLogger)

	newLogger.SendLog("info", fmt.Sprintf("Start running server. Address: %s, db: %s", cfg.RunAddress, cfg.DatabaseURI))

	srv := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			newLogger.SendLog("fatal", fmt.Sprintf("Running server fail %v", zap.Error(err)))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	newLogger.SendLog("info", "Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		newLogger.SendLog("fatal", fmt.Sprintf("Server forced to shutdown %v", zap.Error(err)))
	}

	newLogger.SendLog("info", "Server exiting")
}
