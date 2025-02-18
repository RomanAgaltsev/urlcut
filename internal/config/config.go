// Пакет config формирует конфигурацию приложения.
package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// ErrInitConfigFailed - ошибка инициации конфигурации.
var ErrInitConfigFailed = fmt.Errorf("failed to init config")

// Config - структура конфигурации приложения.
type Config struct {
	ServerPort      string // Адрес HTTP сервера и порт
	BaseURL         string // Базовый адрес сокращенного URL
	FileStoragePath string // Путь к файловому хранилищу
	DatabaseDSN     string // Строка соединения с БД
	SecretKey       string // Секретный ключ авторизации
	EnableHTTPS     bool   // Регулирует включение HTTPS на сервере
	IDlength        int    // Длина идентификатора в сокращенном URL
}

// configBuilder - строитель конфигурации приложения.
type configBuilder struct {
	serverPort      string `env:"SERVER_ADDRESS"`
	baseURL         string `env:"BASE_URL"`
	fileStoragePath string `env:"FILE_STORAGE_PATH"`
	databaseDSN     string `env:"DATABASE_DSN"`
	secretKey       string `env:"SECRET_KEY"`
	enableHTTPS     bool   `env:"ENABLE_HTTPS"`
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
	flag.IntVar(&cb.idLength, "l", cb.idLength, "URL ID default length")
	flag.Parse()

	return nil
}

// setEnvs устанавливает значения конфигурации приложения из переменных окружения.
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
	if dsn != "" {
		cb.secretKey = sk
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
		IDlength:        cb.idLength,
	}
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
