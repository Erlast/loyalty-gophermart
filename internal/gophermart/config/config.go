package config

import (
	"flag"
	"go.uber.org/zap"
	"os"
)

type Config struct {
	RunAddress           string `json:"run_address"`
	Port                 int    `json:"port"`
	DatabaseURI          string `json:"database_uri"`
	AccrualSystemAddress string `json:"accrual_system_address"`
	JWTSecret            string `json:"jwt_secret"`
}

var config Config

func GetConfig() *Config {
	return &config
}

// LoadConfig loads configuration from environment variables and flags.
func LoadConfig(logger *zap.SugaredLogger) Config {
	// Initialize default values from environment variables
	defaultRunAddress := "localhost"
	defaultPort := "8080"
	defaultDatabaseURI := ""
	defaultAccrualSystemAddress := ""
	defaultJWTSecret := "secret"

	if envRunAddress, exists := os.LookupEnv("RUN_ADDRESS"); exists {
		defaultRunAddress = envRunAddress
	}

	if envPort, exists := os.LookupEnv("PORT"); exists {
		defaultPort = envPort
	}

	if envDatabaseURI, exists := os.LookupEnv("DATABASE_URI"); exists {
		defaultDatabaseURI = envDatabaseURI
	}
	if envAccrualSystemAddress, exists := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); exists {
		defaultAccrualSystemAddress = envAccrualSystemAddress
	}
	if envJWTSecret, exists := os.LookupEnv("JWT_SECRET"); exists {
		defaultJWTSecret = envJWTSecret
	}

	// Define command-line flags with default values from environment variables
	runAddress := flag.String("a", defaultRunAddress, "service run address")
	port := flag.String("p", defaultPort, "service port")
	databaseURI := flag.String("d", defaultDatabaseURI, "Database URI")
	accrualSystemAddress := flag.String("r", defaultAccrualSystemAddress, "Accrual system address")
	jwtSecret := flag.String("s", defaultJWTSecret, "JWT secret")

	// Parse the flags
	flag.Parse()

	logger.Infof("database URI: %s, accrualSystemAddres: %s", *databaseURI, *accrualSystemAddress)
	// Ensure that the required parameters are provided
	if *databaseURI == "" || *accrualSystemAddress == "" {
		logger.Fatal("DATABASE URI and ACCRUAL SYSTEM ADDRESS must be provided")
	}

	runAddressWithPort := *runAddress + ":" + *port

	config = Config{
		RunAddress:           runAddressWithPort,
		DatabaseURI:          *databaseURI,
		AccrualSystemAddress: *accrualSystemAddress,
		JWTSecret:            *jwtSecret,
	}

	return config
}
