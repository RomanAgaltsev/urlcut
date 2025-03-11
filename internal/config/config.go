// Пакет config формирует конфигурацию приложения.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// ErrInitConfigFailed - ошибка инициации конфигурации.
var ErrInitConfigFailed = fmt.Errorf("failed to init config")

// Config - структура конфигурации приложения.
type Config struct {
	ServerPort      string `json:"server_address"`    // Адрес HTTP сервера и порт
	BaseURL         string `json:"base_url"`          // Базовый адрес сокращенного URL
	FileStoragePath string `json:"file_storage_path"` // Путь к файловому хранилищу
	DatabaseDSN     string `json:"database_dsn"`      // Строка соединения с БД
	SecretKey       string `json:"secret_key"`        // Секретный ключ авторизации
	EnableHTTPS     bool   `json:"enable_https"`      // Регулирует включение HTTPS на сервере
	TrustedSubnet   string `json:"trusted_subnet"`    // Cтроковое представление бесклассовой адресации доверенной подсети
	ServerGRPCPort  string `json:"server_grpc_port"`  // Порт GRPC сервера
	IDlength        int    `json:"id_length"`         // Длина идентификатора в сокращенном URL
}

// configBuilder - строитель конфигурации приложения.
type configBuilder struct {
	serverPort      string `env:"SERVER_ADDRESS"`
	baseURL         string `env:"BASE_URL"`
	fileStoragePath string `env:"FILE_STORAGE_PATH"`
	databaseDSN     string `env:"DATABASE_DSN"`
	secretKey       string `env:"SECRET_KEY"`
	enableHTTPS     bool   `env:"ENABLE_HTTPS"`
	trustedSubnet   string `env:"TRUSTED_SUBNET"`
	serverGRPCPort  string `env:"SERVER_GRPC_PORT"`
	idLength        int
}

// newConfigBuilder создает нового строителя конфигурации приложения.
func newConfigBuilder() *configBuilder {
	return &configBuilder{}
}

// setDefaults устанавливает значения конфигурации приложения по умолчанию.
func (cb *configBuilder) setDefaults() error {
	cb.serverPort = "localhost:8080"
	cb.baseURL = "http://localhost:8080"
	cb.fileStoragePath = "storage.json"
	cb.databaseDSN = ""
	cb.secretKey = "secret"
	cb.enableHTTPS = false
	cb.trustedSubnet = ""
	cb.serverGRPCPort = ":9090"
	cb.idLength = 8

	return nil
}

// setFlags устанавливает значения конфигурации приложения из параметров командной строки.
func (cb *configBuilder) setFlags() error {
	flag.StringVar(&cb.serverPort, "a", cb.serverPort, "address and port to run server")
	flag.StringVar(&cb.baseURL, "b", cb.baseURL, "basic address of shortened URL")
	flag.StringVar(&cb.fileStoragePath, "f", cb.fileStoragePath, "path to the storage file")
	flag.StringVar(&cb.databaseDSN, "d", cb.databaseDSN, "database connection string")
	flag.StringVar(&cb.secretKey, "k", cb.secretKey, "secret authorization key")
	flag.BoolVar(&cb.enableHTTPS, "s", cb.enableHTTPS, "enable HTTPS on server")
	flag.StringVar(&cb.trustedSubnet, "t", cb.trustedSubnet, "trusted subnet")
	flag.StringVar(&cb.serverGRPCPort, "g", cb.serverGRPCPort, "GRPC server port")
	flag.IntVar(&cb.idLength, "l", cb.idLength, "URL ID default length")
	flag.Parse()

	return nil
}

// setEnvs устанавливает значения конфигурации приложения из переменных окружения.
func (cb *configBuilder) setEnvs() error {
	// Получаем конфигурацию из файла
	configFile := os.Getenv("CONFIG")
	if configFile != "" {
		fromFile, err := configFromFile(configFile)
		if err != nil {
			log.Printf("reading config from file : %s", err.Error())
		} else {
			if fromFile.ServerPort != "" {
				cb.serverPort = fromFile.ServerPort
			}
			if fromFile.BaseURL != "" {
				cb.baseURL = fromFile.BaseURL
			}
			if fromFile.FileStoragePath != "" {
				cb.fileStoragePath = fromFile.FileStoragePath
			}
			if fromFile.DatabaseDSN != "" {
				cb.databaseDSN = fromFile.DatabaseDSN
			}
			if fromFile.SecretKey != "" {
				cb.secretKey = fromFile.SecretKey
			}
			if fromFile.EnableHTTPS {
				cb.enableHTTPS = fromFile.EnableHTTPS
			}
			if fromFile.EnableHTTPS {
				cb.enableHTTPS = fromFile.EnableHTTPS
			}
			if fromFile.TrustedSubnet != "" {
				cb.trustedSubnet = fromFile.TrustedSubnet
			}
			if fromFile.ServerGRPCPort != "" {
				cb.serverGRPCPort = fromFile.ServerGRPCPort
			}
			if fromFile.IDlength != 0 {
				cb.idLength = fromFile.IDlength
			}
		}
	}

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

	eh := os.Getenv("ENABLE_HTTPS")
	if eh != "" {
		enableHTTPS, errConv := strconv.ParseBool(eh)
		if errConv != nil {
			cb.enableHTTPS = false
		} else {
			cb.enableHTTPS = enableHTTPS
		}
	}

	sk := os.Getenv("SECRET_KEY")
	if sk != "" {
		cb.secretKey = sk
	}

	ts := os.Getenv("TRUSTED_SUBNET")
	if ts != "" {
		cb.secretKey = ts
	}

	sgp := os.Getenv("SERVER_GRPC_PORT")
	if sgp != "" {
		cb.serverGRPCPort = sgp
	}

	return nil
}

// build строит конфигурацию приложения.
func (cb *configBuilder) build() *Config {
	return &Config{
		ServerPort:      cb.serverPort,
		BaseURL:         cb.baseURL,
		FileStoragePath: cb.fileStoragePath,
		DatabaseDSN:     cb.databaseDSN,
		SecretKey:       cb.secretKey,
		EnableHTTPS:     cb.enableHTTPS,
		TrustedSubnet:   cb.trustedSubnet,
		ServerGRPCPort:  cb.serverGRPCPort,
		IDlength:        cb.idLength,
	}
}

// configFromFile читает и возвращает конфигурацию приложения из JSON файла.
func configFromFile(fname string) (Config, error) {
	var cfg Config

	data, err := os.ReadFile(fname)
	if err != nil {
		return cfg, err
	}

	if err = json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Get возвращает конфигурацию приложения.
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
			return nil, ErrInitConfigFailed
		}
	}

	return cb.build(), nil
}
