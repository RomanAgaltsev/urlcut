package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name  string
		envs  map[string]string
		flags []string
		cfg   *Config
	}{
		{"envs without flags",
			map[string]string{
				"SERVER_ADDRESS":    "localhost:8081",
				"BASE_URL":          "http://localhost:8081",
				"FILE_STORAGE_PATH": "storage.json",
				"DATABASE_DSN":      "test_dsn",
				"SECRET_KEY":        "secret_key",
				"ENABLE_HTTPS":      "true"},
			[]string{},
			&Config{
				ServerPort:      "localhost:8081",
				BaseURL:         "http://localhost:8081",
				FileStoragePath: "storage.json",
				DatabaseDSN:     "test_dsn",
				SecretKey:       "secret_key",
				EnableHTTPS:     true,
				IDlength:        8,
			},
		},

		{"flags without envs",
			map[string]string{
				"SERVER_ADDRESS":    "",
				"BASE_URL":          "",
				"FILE_STORAGE_PATH": "",
				"DATABASE_DSN":      "",
				"SECRET_KEY":        "",
				"ENABLE_HTTPS":      ""},
			[]string{
				"-a", "localhost:8082",
				"-b", "http://localhost:8082",
				"-f", "storage.json",
				"-d", "test1_dsn",
				"-k", "secret_key1",
				"-s", "true",
				"-l", "8"},
			&Config{
				ServerPort:      "localhost:8080",
				BaseURL:         "http://localhost:8080",
				FileStoragePath: "storage.json",
				DatabaseDSN:     "test1_dsn",
				SecretKey:       "secret_key1",
				EnableHTTPS:     true,
				IDlength:        8,
			},
		},

		{"all envs and all flags",
			map[string]string{
				"SERVER_ADDRESS":    "localhost:8083",
				"BASE_URL":          "http://localhost:8083",
				"FILE_STORAGE_PATH": "storage1.json",
				"DATABASE_DSN":      "test_dsn",
				"SECRET_KEY":        "secret_key",
				"ENABLE_HTTPS":      "true"},
			[]string{
				"-a", "localhost:8084",
				"-b", "http://localhost:8084",
				"-f", "storage2.json",
				"-d", "test1_dsn",
				"-k", "secret_key1",
				"-s", "false",
				"-l", "8"},
			&Config{
				ServerPort:      "localhost:8083",
				BaseURL:         "http://localhost:8083",
				FileStoragePath: "storage1.json",
				DatabaseDSN:     "test_dsn",
				SecretKey:       "secret_key",
				EnableHTTPS:     true,
				IDlength:        8,
			},
		},
		{"envs and flags #1",
			map[string]string{
				"SERVER_ADDRESS":    "localhost:8084",
				"BASE_URL":          "",
				"FILE_STORAGE_PATH": "",
				"DATABASE_DSN":      "",
				"SECRET_KEY":        "",
				"ENABLE_HTTPS":      ""},
			[]string{
				"-b", "http://localhost:8085",
				"-f", "storage2.json",
				"-k", "secret_key1",
				"-s", "true",
				"-d", "test1_dsn"},
			&Config{
				ServerPort:      "localhost:8084",
				BaseURL:         "http://localhost:8080",
				FileStoragePath: "storage2.json",
				DatabaseDSN:     "",
				SecretKey:       "secret_key1",
				EnableHTTPS:     true,
				IDlength:        8,
			},
		},
		{"envs and flags #2",
			map[string]string{
				"SERVER_ADDRESS":    "",
				"BASE_URL":          "http://localhost:8086",
				"FILE_STORAGE_PATH": "storage1.json",
				"DATABASE_DSN":      "test_dsn",
				"SECRET_KEY":        "secret_key",
				"ENABLE_HTTPS":      "true"},
			[]string{
				"-a", "localhost:8087",
				"-l", "8"},
			&Config{
				ServerPort:      "localhost:8080",
				BaseURL:         "http://localhost:8086",
				FileStoragePath: "storage1.json",
				DatabaseDSN:     "test_dsn",
				SecretKey:       "secret_key",
				EnableHTTPS:     true,
				IDlength:        8,
			},
		},
		{"envs and flags #3",
			map[string]string{
				"SERVER_ADDRESS":    "localhost:8088",
				"BASE_URL":          "http://localhost:8088",
				"FILE_STORAGE_PATH": "storage1.json",
				"DATABASE_DSN":      "test_dsn",
				"SECRET_KEY":        "secret_key",
				"ENABLE_HTTPS":      "true"},
			[]string{
				"-a", "localhost:8089"},
			&Config{
				ServerPort:      "localhost:8088",
				BaseURL:         "http://localhost:8088",
				FileStoragePath: "storage1.json",
				DatabaseDSN:     "test_dsn",
				SecretKey:       "secret_key",
				EnableHTTPS:     true,
				IDlength:        8,
			},
		},
		{"envs and flags #4",
			map[string]string{
				"SERVER_ADDRESS": "",
				"BASE_URL":       "http://localhost:8090"},
			[]string{
				"-a", "localhost:8091",
				"-b", "http://localhost:8091",
				"-f", "storage2.json",
				"-d", "test1_dsn",
				"-k", "secret_key1",
				"-s", "true",
				"-l", "8"},
			&Config{
				ServerPort:      "localhost:8080",
				BaseURL:         "http://localhost:8090",
				FileStoragePath: "storage2.json",
				DatabaseDSN:     "test1_dsn",
				SecretKey:       "secret_key1",
				EnableHTTPS:     true,
				IDlength:        8,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defArgs := os.Args
			defCL := flag.CommandLine

			defer func() {
				os.Args = defArgs
				flag.CommandLine = defCL
			}()

			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			t.Setenv("CONFIG", "config.json")

			for k, v := range test.envs {
				t.Setenv(k, v)
			}

			os.Args = append([]string{"cmd"}, test.flags...)

			cfg, err := Get()
			require.NoError(t, err)
			assert.Equal(t, test.cfg, cfg)
		})
	}
}
