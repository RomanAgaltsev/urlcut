package main

import (
    "github.com/RomanAgaltsev/urlcut/internal/app"
    "github.com/RomanAgaltsev/urlcut/internal/config"
    "log"
)

func main() {
    // Получаем конфигурацию
    cfg, err := config.Get()
    // Проверяем на возможные ошибки
    if err != nil {
        // Есть ошибка получения конфигурации
        log.Printf("getting config failed: %v", err)
        return
    }
    // Запускаем приложение
    app.Run(cfg)
}
