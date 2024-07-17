package main

import (
	"context"
	"net/http"

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

	go components.OrderProcessing(ctx, store, newLogger)

	r := routes.NewAccrualRouter(ctx, store, newLogger)

	newLogger.Infof("Start running server. Address: %s, db: %s", cfg.RunAddress, cfg.DatabaseURI)
	err = http.ListenAndServe(cfg.RunAddress, r)

	if err != nil {
		newLogger.Fatalf("Running server fail %s", err)
	}
}
