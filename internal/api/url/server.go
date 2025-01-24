package url

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	"github.com/RomanAgaltsev/urlcut/internal/api/middleware"
	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
)

// ErrInitServerFailed ошибка инициализации HTTP сервера.
var ErrInitServerFailed = fmt.Errorf("failed to init HTTP server")

// NewServer создает новый HTTP сервер с установкой обработчиков и роутера.
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
	// -- идентификатор пользователя не требуется - выдаем при отсутствии
	tokenAuth := jwtauth.New("HS256", []byte(cfg.SecretKey), nil)
	router.Group(func(r chi.Router) {
		// Миддлваре, проверяющая и выдающая токен
		r.Use(middleware.WithAuth(tokenAuth))

		r.Post("/", handlers.Shorten)
		r.Post("/api/shorten", handlers.ShortenAPI)
		r.Post("/api/shorten/batch", handlers.ShortenAPIBatch)
		r.Get("/{id}", handlers.Expand)
		r.Get("/ping", handlers.Ping)
	})
	// -- идентификатор требуется
	router.Group(func(r chi.Router) {
		// Миддвале, проверяющая наличие идентификтара
		r.Use(middleware.WithID(tokenAuth))

		r.Get("/api/user/urls", handlers.UserUrls)
		r.Delete("/api/user/urls", handlers.UserUrlsDelete)
	})

	return &http.Server{
		Addr:    cfg.ServerPort,
		Handler: router,
	}, nil
}
