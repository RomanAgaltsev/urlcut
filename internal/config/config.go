package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

// Cfg - структура для хранения конфигурации
type Cfg struct {
	ServerAddr string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`  // Адрес HTTP-сервера
	BasicAddr  string `env:"BASE_URL" envDefault:"http://localhost:8080"` // Базовый адрес результирующего сокращённого URL
	IDlength   int    // Длина идентификатора сокращенного URL
}

// Config - хранит значения конфигурации
var Config Cfg

// ParseFlags - выполняет парсинг флагов
func ParseFlags() {
	// Создаем структуру
	Config = Cfg{}

	if err := env.Parse(&Config); err != nil {
		fmt.Printf("Error parsing environment variables: %+v\n", err)
	}

	flag.StringVar(&Config.ServerAddr, "a", Config.ServerAddr, "address and port to run server")
	flag.StringVar(&Config.BasicAddr, "b", Config.BasicAddr, "basic address of shortened URL")
	flag.IntVar(&Config.IDlength, "l", 8, "URL ID default length")

	// Парсим флаги
	flag.Parse()
}
