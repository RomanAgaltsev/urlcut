package config

import (
    "flag"
    "fmt"
    "github.com/caarlos0/env/v6"
)

// Config - структура для хранения конфигурации
type Config struct {
    ServerPort string // Адрес сервера и его порт
    BaseURL    string // Базовый URL для сокращенного URL
    IDlength   int    // Длина идентификатора сокращенного URL
}

// configBuilder - строитель конфигурации, поля дублируют Config
type configBuilder struct {
    serverPort string `env:"SERVER_ADDRESS"`
    baseURL    string `env:"BASE_URL"`
    idLength   int
}

// newConfigBuilder - возвращает нового строителя конфигурации
func newConfigBuilder() *configBuilder {
    return &configBuilder{}
}

// setDefaults - устанавливает параметры конфигурации по умолчанию
func (cb *configBuilder) setDefaults() error {
    cb.serverPort = "localhost:8080"
    cb.baseURL = "http://localhost:8080"
    cb.idLength = 8
    return nil
}

// setFlags - устанавливает параметры конфигурации из параметров командной строки
func (cb *configBuilder) setFlags() error {
    flag.StringVar(&cb.serverPort, "a", cb.serverPort, "address and port to run server")
    flag.StringVar(&cb.baseURL, "b", cb.baseURL, "basic address of shortened URL")
    flag.IntVar(&cb.idLength, "l", cb.idLength, "URL ID default length")
    flag.Parse()
    return nil
}

// setEnvs - устанавливает параметры конфигурации из переменных окружения
func (cb *configBuilder) setEnvs() error {
    if err := env.Parse(cb); err != nil {
        return fmt.Errorf("Error parsing environment variables: %+v\n", err)
    }
    return nil
}

// getConfig - возвращает заполненную структуру конфигурации
func (cb *configBuilder) getConfig() *Config {
    return &Config{
        ServerPort: cb.serverPort,
        BaseURL:    cb.baseURL,
        IDlength:   cb.idLength,
    }
}

// Get - возвращает структуру конфигурации, заполненную по параметрам командной строки или переменным окружения
func Get() (*Config, error) {
    cb := newConfigBuilder()

    err := cb.setDefaults()
    if err != nil {
        return nil, err
    }

    err = cb.setFlags()
    if err != nil {
        return nil, err
    }

    err = cb.setEnvs()
    if err != nil {
        return nil, err
    }

    return cb.getConfig(), nil
}
