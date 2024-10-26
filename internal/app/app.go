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

	"github.com/RomanAgaltsev/urlcut/internal/api/middleware"
	apiurl "github.com/RomanAgaltsev/urlcut/internal/api/url"
	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/logger"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
	repositoryurl "github.com/RomanAgaltsev/urlcut/internal/repository/url"
	services "github.com/RomanAgaltsev/urlcut/internal/service"
	servicesurl "github.com/RomanAgaltsev/urlcut/internal/service/url"

	"github.com/go-chi/chi/v5"
)

var (
	ErrInitConfigFailed  = fmt.Errorf("failed to init config")
	ErrInitServiceFailed = fmt.Errorf("failed to init service")
	ErrInitServerFailed  = fmt.Errorf("failed to init HTTP server")
)

type App struct {
	cfg     *config.Config
	repo    repository.URLRepository
	service services.URLService
	server  *http.Server
}

func New() (*App, error) {
	app := &App{}

	err := app.init()
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) init() error {
	appInits := []func() error{
		a.initConfig,
		a.initLogger,
		a.initRepository,
		a.initService,
		a.initHTTPServer,
	}

	for _, appInit := range appInits {
		err := appInit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig() error {
	cfg, err := config.Get()
	if err != nil {
		return ErrInitConfigFailed
	}
	a.cfg = cfg

	return nil
}

func (a *App) initLogger() error {
	err := logger.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initRepository() error {
	inMemoRepo := repositoryurl.New(a.cfg.FileStoragePath)

	if err := inMemoRepo.RestoreState(); err != nil {
		return err
	}
	a.repo = inMemoRepo

	return nil
}

func (a *App) initService() error {
	if a.cfg.BaseURL == "" || a.cfg.IDlength == 0 {
		return ErrInitServiceFailed
	}

	a.service = servicesurl.New(a.repo, a.cfg.BaseURL, a.cfg.IDlength)

	return nil
}

func (a *App) initHTTPServer() error {
	if a.cfg.ServerPort == "" {
		return ErrInitServerFailed
	}
	handlers := apiurl.New(a.service)

	router := chi.NewRouter()
	router.Use(middleware.WithLogging)
	router.Use(middleware.WithGzip)
	router.Post("/", handlers.Shorten)
	router.Post("/api/shorten", handlers.ShortenAPI)
	router.Get("/{id}", handlers.Expand)

	a.server = &http.Server{
		Addr:    a.cfg.ServerPort,
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

		if err := a.repo.SaveState(); err != nil {
			slog.Error(
				"failed to save url storage to file",
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
