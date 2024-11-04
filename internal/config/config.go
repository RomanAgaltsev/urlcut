package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerPort      string
	BaseURL         string
	FileStoragePath string
	DatabaseDSN     string
	IDlength        int
}

type configBuilder struct {
	serverPort      string `env:"SERVER_ADDRESS"`
	baseURL         string `env:"BASE_URL"`
	fileStoragePath string `env:"FILE_STORAGE_PATH"`
	databaseDSN     string `env:"DATABASE_DSN"`
	idLength        int
}

func newConfigBuilder() *configBuilder {
	return &configBuilder{}
}

func (cb *configBuilder) setDefaults() error {
	cb.serverPort = "localhost:8080"
	cb.baseURL = "http://localhost:8080"
	cb.fileStoragePath = "storage.json"
	cb.databaseDSN = ""
	cb.idLength = 8

	return nil
}

func (cb *configBuilder) setFlags() error {
	flag.StringVar(&cb.serverPort, "a", cb.serverPort, "address and port to run server")
	flag.StringVar(&cb.baseURL, "b", cb.baseURL, "basic address of shortened URL")
	flag.StringVar(&cb.fileStoragePath, "f", cb.fileStoragePath, "path to the storage file")
	flag.StringVar(&cb.databaseDSN, "d", cb.databaseDSN, "database connection string")
	flag.IntVar(&cb.idLength, "l", cb.idLength, "URL ID default length")
	flag.Parse()

	return nil
}

func (cb *configBuilder) setEnvs() error {
	sp := os.Getenv("SERVER_ADDRESS")
	if sp != "" {
		cb.serverPort = sp
	}

	bu := os.Getenv("BASE_URL")
	if bu != "" {
		cb.baseURL = bu
	}

	fsp := os.Getenv("FILE_STORAGE_PATH")
	if fsp != "" {
		cb.fileStoragePath = fsp
	}

	dsn := os.Getenv("DATABASE_DSN")
	if dsn != "" {
		cb.databaseDSN = dsn
	}

	return nil
}

func (cb *configBuilder) build() *Config {
	return &Config{
		ServerPort:      cb.serverPort,
		BaseURL:         cb.baseURL,
		FileStoragePath: cb.fileStoragePath,
		DatabaseDSN:     cb.databaseDSN,
		IDlength:        cb.idLength,
	}
}

func Get() (*Config, error) {
	cb := newConfigBuilder()

	confSets := []func() error{
		cb.setDefaults,
		cb.setFlags,
		cb.setEnvs,
	}

	for _, confSet := range confSets {
		err := confSet()
		if err != nil {
			return nil, err
		}
	}

	return cb.build(), nil
}
