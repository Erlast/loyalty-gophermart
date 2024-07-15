package config

import (
	"flag"

	"os"
	"testing"

	"github.com/Erlast/loyalty-gophermart.git/pkg/zaplog"
	"github.com/stretchr/testify/assert"
)

func TestParseAccrualFlags(t *testing.T) {
	zaplog.InitLogger()
	os.Args = []string{"cmd", "-a", "127.0.0.1:9090", "-d", "postgres://user:pass@localhost/db"}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	config := ParseFlags()

	assert.Equal(t, "127.0.0.1:9090", config.RunAddress)
	assert.Equal(t, "postgres://user:pass@localhost/db", config.DatabaseURI)

	os.Args = []string{"cmd"}
	t.Setenv("RUN_ADDRESS", "192.168.1.1:8080")
	t.Setenv("DATABASE_URI", "postgres://envuser:envpass@localhost/envdb")

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	config = ParseFlags()

	assert.Equal(t, "192.168.1.1:8080", config.RunAddress)
	assert.Equal(t, "postgres://envuser:envpass@localhost/envdb", config.DatabaseURI)
}
