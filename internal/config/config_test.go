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
				"SERVER_ADDRESS": "localhost:8081",
				"BASE_URL":       "http://localhost:8081"},
			[]string{},
			&Config{
				ServerPort: "localhost:8081",
				BaseURL:    "http://localhost:8081",
				IDlength:   8,
			},
		},

		{"flags without envs",
			map[string]string{
				"SERVER_ADDRESS": "",
				"BASE_URL":       ""},
			[]string{
				"-a", "localhost:8082",
				"-b", "http://localhost:8082",
				"-l", "9"},
			&Config{
				ServerPort: "localhost:8082",
				BaseURL:    "http://localhost:8082",
				IDlength:   9,
			},
		},

		{"all envs and all flags",
			map[string]string{
				"SERVER_ADDRESS": "localhost:8083",
				"BASE_URL":       "http://localhost:8083"},
			[]string{
				"-a", "localhost:8084",
				"-b", "http://localhost:8084",
				"-l", "10"},
			&Config{
				ServerPort: "localhost:8083",
				BaseURL:    "http://localhost:8083",
				IDlength:   10,
			},
		},
		{"envs and flags #1",
			map[string]string{
				"SERVER_ADDRESS": "localhost:8084",
				"BASE_URL":       ""},
			[]string{
				"-b", "http://localhost:8085"},
			&Config{
				ServerPort: "localhost:8084",
				BaseURL:    "http://localhost:8085",
				IDlength:   8,
			},
		},
		{"envs and flags #2",
			map[string]string{
				"SERVER_ADDRESS": "",
				"BASE_URL":       "http://localhost:8086"},
			[]string{
				"-a", "localhost:8087",
				"-l", "12"},
			&Config{
				ServerPort: "localhost:8087",
				BaseURL:    "http://localhost:8086",
				IDlength:   12,
			},
		},
		{"envs and flags #3",
			map[string]string{
				"SERVER_ADDRESS": "localhost:8088",
				"BASE_URL":       "http://localhost:8088"},
			[]string{
				"-a", "localhost:8089"},
			&Config{
				ServerPort: "localhost:8088",
				BaseURL:    "http://localhost:8088",
				IDlength:   8,
			},
		},
		{"envs and flags #4",
			map[string]string{
				"SERVER_ADDRESS": "",
				"BASE_URL":       "http://localhost:8090"},
			[]string{
				"-a", "localhost:8091",
				"-b", "http://localhost:8091",
				"-l", "12"},
			&Config{
				ServerPort: "localhost:8091",
				BaseURL:    "http://localhost:8090",
				IDlength:   12,
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
