package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	apiurl "github.com/RomanAgaltsev/urlcut/internal/api/url"
	repository "github.com/RomanAgaltsev/urlcut/internal/repository"
	repositoryurl "github.com/RomanAgaltsev/urlcut/internal/repository/url"
	services "github.com/RomanAgaltsev/urlcut/internal/service"
	servicesurl "github.com/RomanAgaltsev/urlcut/internal/service/url"

	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/logger"

	"github.com/go-chi/chi/v5"
)

type App struct {
	repo    repository.URLRepository
	service services.URLService
	server  *http.Server
}

func New(_ context.Context) (*App, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config : %v", err)
	}
	err = logger.Initialize()
	if err != nil {
		return nil, err
	}
	app := &App{}
	err = app.getRepository()
	if err != nil {
		return nil, err
	}
	err = app.getService(cfg.BaseURL, cfg.IDlength)
	if err != nil {
		return nil, err
	}
	err = app.getHTTPServer(cfg.ServerPort)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) getRepository() error {
	a.repo = repositoryurl.New()
	return nil
}

func (a *App) getService(baseURL string, idLength int) error {
	a.service = servicesurl.NewShortener(a.repo, baseURL, idLength)
	return nil
}

func (a *App) getHTTPServer(serverPort string) error {
	handlers := apiurl.NewHandlers(a.service)
	router := chi.NewRouter()
	router.Post("/", apiurl.WithLogging(handlers.ShortenURL))
	router.Get("/{id}", apiurl.WithLogging(handlers.ExpandURL))
	a.server = &http.Server{
		Addr:    serverPort,
		Handler: router,
	}
	return nil
}

func (a *App) Run() error {
	return a.runShortenerApp()
}

func (a *App) runShortenerApp() error {
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		slog.Info("shutting down HTTP server")

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if err := a.server.Shutdown(ctx); err != nil {
			slog.Error(
				"HTTP server shutdown error",
				slog.String("error", err.Error()),
			)
		}
		close(done)
	}()

	slog.Info(
		"starting HTTP server",
		"addr", a.server.Addr,
	)
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error(
			"HTTP server error",
			slog.String("error", err.Error()),
		)
		return err
	}

	<-done
	slog.Info("HTTP server stopped")
	return nil
}
