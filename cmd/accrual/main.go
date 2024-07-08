package main

import (
	"context"
	"github.com/Erlast/loyalty-gophermart.git/internal/config"
	"github.com/Erlast/loyalty-gophermart.git/internal/log"
	"github.com/Erlast/loyalty-gophermart.git/internal/routes"
	"github.com/Erlast/loyalty-gophermart.git/internal/storage"
	"net/http"
)

func main() {
	ctx := context.Background()
	cfg := config.ParseAccrualFlags()

	logger, err := log.NewLogger("info")

	store, err := storage.NewAccrualStorage(ctx, cfg, logger)
	if err != nil {
		logger.Fatalf("Unable to create storage %v: ", err)
	}

	r := routes.NewAccrualRouter(ctx, cfg, store, logger)

	logger.Infof("Start running server. Address: %s, db: %s", cfg.RunAddress, cfg.DatabaseURI)
	err = http.ListenAndServe(cfg.RunAddress, r)

	if err != nil {
		logger.Fatalf("Running server fail %s", err)
	}
}
