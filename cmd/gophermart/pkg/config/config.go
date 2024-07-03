package config

import (
	"flag"
	"gofermart/cmd/gophermart/pkg/zaplog"
	"os"
)

type Config struct {
	RunAddress           string `json:"run_address"`
	DatabaseURI          string `json:"database_uri"`
	AccrualSystemAddress string `json:"accrual_system_address"`
}

// LoadConfig loads configuration from environment variables and flags
func LoadConfig() Config {
	// Initialize default values from environment variables
	defaultRunAddress := "localhost:8080"
	defaultDatabaseURI := ""
	defaultAccrualSystemAddress := ""

	if envRunAddress, exists := os.LookupEnv("RUN_ADDRESS"); exists {
		defaultRunAddress = envRunAddress
	}
	if envDatabaseURI, exists := os.LookupEnv("DATABASE_URI"); exists {
		defaultDatabaseURI = envDatabaseURI
	}
	if envAccrualSystemAddress, exists := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); exists {
		defaultAccrualSystemAddress = envAccrualSystemAddress
	}

	// Define command-line flags with default values from environment variables
	runAddress := flag.String("a", defaultRunAddress, "Service run address")
	databaseURI := flag.String("d", defaultDatabaseURI, "Database URI")
	accrualSystemAddress := flag.String("r", defaultAccrualSystemAddress, "Accrual system address")

	// Parse the flags
	flag.Parse()

	// Ensure that the required parameters are provided
	if *databaseURI == "" || *accrualSystemAddress == "" {
		zaplog.Logger.Fatalf("Both DATABASE_URI and ACCRUAL_SYSTEM_ADDRESS must be provided either as flags or environment variables")
	}

	return Config{
		RunAddress:           *runAddress,
		DatabaseURI:          *databaseURI,
		AccrualSystemAddress: *accrualSystemAddress,
	}
}
