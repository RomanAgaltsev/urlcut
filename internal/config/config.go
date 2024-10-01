package config

import (
	"flag"
	"os"
)

// Cfg - структура для хранения конфигурации
type Cfg struct {
	ServerAddr string // Адрес HTTP-сервера
	BasicAddr  string // Базовый адрес результирующего сокращённого URL
	IDlength   int    // Длина идентификатора сокращенного URL
}

// Config - хранит значения конфигурации
var Config Cfg

// ParseFlags - выполняет парсинг флагов
func ParseFlags() {
	// Создаем структуру
	Config = Cfg{}

	// Устанавливам соответствие полей структуры и переменных окружения/флагов
	if serverAddr := os.Getenv("SERVER_ADDRESS"); serverAddr != "" {
		Config.ServerAddr = serverAddr
	} else {
		flag.StringVar(&Config.ServerAddr, "a", "localhost:8080", "address and port to run server")
	}
	if basicAddr := os.Getenv("BASE_URL"); basicAddr != "" {
		Config.BasicAddr = basicAddr
	} else {
		flag.StringVar(&Config.BasicAddr, "b", "http://localhost:8080", "basic address of shortened URL")
	}
	flag.IntVar(&Config.IDlength, "l", 8, "URL ID default length")

	// Парсим флаги
	flag.Parse()
}
