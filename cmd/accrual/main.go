package main

import (
	"context"
	"log"

	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/config"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/logger"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/routes"
	"github.com/Erlast/loyalty-gophermart.git/internal/accrual/storage"

	"net/http"
)

func main() {
	ctx := context.Background()
	cfg := config.ParseAccrualFlags()

	newLogger, err := logger.NewLogger("info")
	if err != nil {
		log.Fatal("Running logger error")
	}

	store, err := storage.NewAccrualStorage(ctx, cfg, newLogger)
	if err != nil {
		newLogger.Fatalf("Unable to create storage %v: ", err)
	}

	r := routes.NewAccrualRouter(ctx, store, newLogger)

	newLogger.Infof("Start running server. Address: %s, db: %s", cfg.RunAddress, cfg.DatabaseURI)
	err = http.ListenAndServe(cfg.RunAddress, r)

	if err != nil {
		newLogger.Fatalf("Running server fail %s", err)
	}
}
