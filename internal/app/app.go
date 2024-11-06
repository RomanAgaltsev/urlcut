package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/RomanAgaltsev/urlcut/internal/api/url"
	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/logger"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
	"github.com/RomanAgaltsev/urlcut/internal/services"
)

type App struct {
	config    *config.Config
	server    *http.Server
	shortener interfaces.Service
}

func New() (*App, error) {
	app := &App{}

	err := app.initConfig()
	if err != nil {
		return nil, err
	}

	err = app.initLogger()
	if err != nil {
		return nil, err
	}

	err = app.initShortener()
	if err != nil {
		return nil, err
	}

	err = app.initHTTPServer()
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) initConfig() error {
	cfg, err := config.Get()
	if err != nil {
		return err
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

func (a *App) initShortener() error {
	repo, err := repository.New(a.config.DatabaseDSN, a.config.FileStoragePath)
	if err != nil {
		return nil
	}

	shortener, err := services.NewShortener(repo, a.config.BaseURL, a.config.IDlength)
	if err != nil {
		return err
	}

	a.shortener = shortener

	return nil
}

func (a *App) initHTTPServer() error {
	server, err := url.NewServer(a.shortener, a.config.ServerPort)
	if err != nil {
		return err
	}
	a.server = server

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

		if err := a.shortener.Close(); err != nil {
			slog.Error(
				"failed to close shortener service",
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
