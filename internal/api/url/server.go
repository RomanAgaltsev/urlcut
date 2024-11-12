package url

import (
	"fmt"

	"net/http"

	"github.com/RomanAgaltsev/urlcut/internal/api/middleware"
	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"

	"github.com/go-chi/chi/v5"
)

// ErrInitServerFailed ошибка инициализации HTTP сервера
var ErrInitServerFailed = fmt.Errorf("failed to init HTTP server")

// NewServer создает новый HTTP сервер с установкой обработчиков и роутера
func NewServer(shortener interfaces.Service, cfg *config.Config) (*http.Server, error) {
	// Если не передали, то ошибка - по умолчанию в конфиге должен быть
	if cfg.ServerPort == "" {
		return nil, ErrInitServerFailed
	}

	// Создаем обработчики
	handlers := NewHandlers(shortener, cfg)

	// Создаем роутер
	router := chi.NewRouter()
	// Включаем миддлаваре
	router.Use(middleware.WithLogging)
	router.Use(middleware.WithGzip)
	// Настраиваем роутинг
	router.Route("/", func(r chi.Router) {
		r.Post("/", handlers.Shorten)
		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", handlers.ShortenAPI)
			r.Post("/shorten/batch", handlers.ShortenAPIBatch)
		})
	})
	router.Get("/{id}", handlers.Expand)
	router.Get("/ping", handlers.Ping)

	return &http.Server{
		Addr:    cfg.ServerPort,
		Handler: router,
	}, nil
}
