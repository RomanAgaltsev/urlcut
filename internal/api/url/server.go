package url

import (
	"fmt"
	"net/http"

	"github.com/RomanAgaltsev/urlcut/internal/api/middleware"
	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
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
	// -- авторизаци не требуется
	router.Group(func(r chi.Router) {
		r.Post("/", handlers.Shorten)
		r.Post("/api/shorten", handlers.ShortenAPI)
		r.Post("/api/shorten/batch", handlers.ShortenAPIBatch)
	})
	// -- авторизация требуется
	tokenAuth := jwtauth.New("HS256", []byte(cfg.SecretKey), nil)
	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(middleware.WithAuth(tokenAuth))

		r.Get("/api/user/urls", handlers.UserUrls)
	})

	return &http.Server{
		Addr:    cfg.ServerPort,
		Handler: router,
	}, nil
}
