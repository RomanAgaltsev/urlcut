package url

import (
	"fmt"
	"net/http"

	"github.com/RomanAgaltsev/urlcut/internal/api/middleware"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"

	"github.com/go-chi/chi/v5"
)

var ErrInitServerFailed = fmt.Errorf("failed to init HTTP server")

func NewServer(shortener interfaces.Service, serverPort string) (*http.Server, error) {
	if serverPort == "" {
		return nil, ErrInitServerFailed
	}
	handlers := NewHandlers(shortener)

	router := chi.NewRouter()
	router.Use(middleware.WithLogging)
	router.Use(middleware.WithGzip)

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
		Addr:    serverPort,
		Handler: router,
	}, nil
}
