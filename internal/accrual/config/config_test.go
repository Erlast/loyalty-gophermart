package config

import (
	"flag"

	"go.uber.org/zap/zaptest"

	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAccrualFlags(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	tests := []struct {
		name         string
		args         []string
		envVars      map[string]string
		expectedAddr string
		expectedURI  string
	}{
		{
			name: "Default values",
			args: []string{},
			envVars: map[string]string{
				"RUN_ADDRESS": "localhost:8080",
			},
			expectedAddr: defaultRunAddress,
			expectedURI:  "",
		},
		{
			name:         "Flag values",
			args:         []string{"-a", "127.0.0.1:9000", "-d", "postgres://user:password@localhost:5432/db"},
			envVars:      map[string]string{},
			expectedAddr: "127.0.0.1:9000",
			expectedURI:  "postgres://user:password@localhost:5432/db",
		},
		{
			name: "Environment values",
			args: []string{},
			envVars: map[string]string{
				"RUN_ADDRESS":  "0.0.0.0:8000",
				"DATABASE_URI": "postgres://user:password@localhost:5432/testdb",
			},
			expectedAddr: "0.0.0.0:8000",
			expectedURI:  "postgres://user:password@localhost:5432/testdb",
		},
		{
			name: "Flag overrides environment",
			args: []string{"-a", "127.0.0.1:9000"},
			envVars: map[string]string{
				"RUN_ADDRESS":  "0.0.0.0:8000",
				"DATABASE_URI": "postgres://user:password@localhost:5432/testdb",
			},
			expectedAddr: "0.0.0.0:8000",
			expectedURI:  "postgres://user:password@localhost:5432/testdb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			var cleanup []func()
			for k, v := range tt.envVars {
				if err := os.Setenv(k, v); err != nil {
					t.Fatalf("failed to set env var: %v", err)
				}
				cleanup = append(cleanup, func(key string) func() {
					return func() {
						if err := os.Unsetenv(key); err != nil {
							t.Fatalf("failed to unset env var: %v", err)
						}
					}
				}(k))
			}

			defer func() {
				for _, fn := range cleanup {
					fn()
				}
			}()

			os.Args = append([]string{os.Args[0]}, tt.args...)

			config := ParseFlags(logger)

			assert.Equal(t, tt.expectedAddr, config.RunAddress)
			assert.Equal(t, tt.expectedURI, config.DatabaseURI)
		})
	}
}
