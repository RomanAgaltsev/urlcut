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
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/logger"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
	"github.com/RomanAgaltsev/urlcut/internal/services"

	"github.com/go-chi/chi/v5"
)

var (
	ErrInitConfigFailed  = fmt.Errorf("failed to init config")
	ErrInitServiceFailed = fmt.Errorf("failed to init service")
	ErrInitServerFailed  = fmt.Errorf("failed to init HTTP server")
)

type App struct {
	config *config.Config
	server *http.Server

	shortener  interfaces.URLShortExpander
	repository interfaces.URLStoreGetter
	stater     interfaces.StateSetGetter
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
		a.initShortener,
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
	a.config = cfg

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
	inMemoryRepository := repository.NewInMemoryRepository()

	saver := services.NewStateSaver(a.config.FileStoragePath)
	state, err := saver.RestoreState()
	if err == nil {
		if err := inMemoryRepository.SetState(state); err != nil {
			slog.Error(
				"failed to restore url storage from file",
				slog.String("error", err.Error()),
			)
		}
	}
	a.repository = inMemoryRepository
	a.stater = inMemoryRepository

	return nil
}

func (a *App) initShortener() error {
	if a.config.BaseURL == "" || a.config.IDlength == 0 {
		return ErrInitServiceFailed
	}

	a.shortener = services.NewShortener(a.repository, a.config.BaseURL, a.config.IDlength)

	return nil
}

func (a *App) initHTTPServer() error {
	if a.config.ServerPort == "" {
		return ErrInitServerFailed
	}
	handlers := apiurl.New(a.shortener)

	router := chi.NewRouter()
	router.Use(middleware.WithLogging)
	router.Use(middleware.WithGzip)
	router.Post("/", handlers.Shorten)
	router.Post("/api/shorten", handlers.ShortenAPI)
	router.Get("/{id}", handlers.Expand)
	router.Get("/ping", handlers.Ping)

	a.server = &http.Server{
		Addr:    a.config.ServerPort,
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

		saver := services.NewStateSaver(a.config.FileStoragePath)
		if err := saver.SaveState(a.stater.GetState()); err != nil {
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
