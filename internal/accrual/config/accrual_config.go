package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v11"
)

type Cfg struct {
	RunAddress  string
	DatabaseURI string
}

type envCfg struct {
	RunAddress  string `env:"RUN_ADDRESS"`
	DatabaseURI string `env:"DATABASE_URI"`
}

const defaultRunAddress = ":8080"

func ParseAccrualFlags() *Cfg {
	config := &Cfg{
		RunAddress:  defaultRunAddress,
		DatabaseURI: "",
	}

	flag.StringVar(&config.RunAddress, "a", config.RunAddress, "port to run server")
	flag.StringVar(&config.DatabaseURI, "d", config.DatabaseURI, "database URI")

	flag.Parse()
	cfg := envCfg{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("can't parse")
	}

	if len(cfg.RunAddress) != 0 {
		config.RunAddress = cfg.RunAddress
	}

	return config
}
