package config

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/Erlast/loyalty-gophermart.git/pkg/opensearch"
)

func TestParseAccrualFlags(t *testing.T) {
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
	os.Args = []string{"cmd", "-a", "127.0.0.1:9090", "-d", "postgres://user:pass@localhost/db"}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	config := ParseFlags(newLogger)

	assert.Equal(t, "127.0.0.1:9090", config.RunAddress)
	assert.Equal(t, "postgres://user:pass@localhost/db", config.DatabaseURI)

	os.Args = []string{"cmd"}
	t.Setenv("RUN_ADDRESS", "192.168.1.1:8080")
	t.Setenv("DATABASE_URI", "postgres://envuser:envpass@localhost/envdb")

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	config = ParseFlags(newLogger)

	assert.Equal(t, "192.168.1.1:8080", config.RunAddress)
	assert.Equal(t, "postgres://envuser:envpass@localhost/envdb", config.DatabaseURI)
}
