package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Erlast/loyalty-gophermart.git/pkg/zaplog"
	"go.uber.org/zap"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/components"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/config"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/routes"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newLogger := zaplog.InitLogger()
	cfg := config.ParseFlags(newLogger)

	store, err := storage.NewAccrualStorage(ctx, cfg)
	if err != nil {
		newLogger.Fatalf("Unable to create storage %v: ", err)
	}

	defer store.DB.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		components.OrderProcessing(ctx, store, newLogger)
	}()

	r := routes.NewAccrualRouter(store, newLogger)

	newLogger.Infof("Start running server. Address: %s, db: %s", cfg.RunAddress, cfg.DatabaseURI)

	srv := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: r,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			newLogger.Fatal("Running server fail", zap.Error(err))
		}
	}()

	go func() {
		<-signalChan
		newLogger.Info("Shutting down server...")

		cancel()

		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		if err := srv.Shutdown(ctxShutdown); err != nil {
			newLogger.Fatal("Server forced to shutdown", zap.Error(err))
		}

		newLogger.Info("Server exited gracefully")
	}()

	wg.Wait()
	newLogger.Info("Server exiting")
}
