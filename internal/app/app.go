package app

import (
    "github.com/RomanAgaltsev/urlcut/internal/api"
    "github.com/RomanAgaltsev/urlcut/internal/config"
    "github.com/RomanAgaltsev/urlcut/internal/repository"
    "github.com/RomanAgaltsev/urlcut/internal/service"
)

// Run - запускает приложение
func Run(cfg *config.Config) {
    // Создаем новое хранилище-мапу
    mapRepository := repository.NewMap()
    // Создаем новый сервис сокращения URL
    shortenerService := service.NewShortener(mapRepository, cfg)
    // Создаем новый HTTP-сервер
    handler := api.NewHandler(shortenerService, cfg)
    // Запускаем сервер
    handler.Run()
}
