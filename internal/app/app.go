// Пакет app обеспечивает создание и запуск всего приложения.
package app

import (
	"context"
	"errors"
	"github.com/RomanAgaltsev/urlcut/internal/pkg/cert"
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

// App является структурой всего приложения.
type App struct {
	cfg       *config.Config     // конфигурация приложения
	server    *http.Server       // http-сервер
	shortener interfaces.Service // сервис сокращателя ссылок
}

// New создает новое приложение.
// При создании приложения инициализируются:
// - конфигурация
// - логер
// - сервис сокращателя ссылок
// - http-сервер
func New() (*App, error) {
	app := &App{}

	// Инициализация конфигурации
	err := app.initConfig()
	if err != nil {
		return nil, err
	}

	// Инициализация логера
	err = app.initLogger()
	if err != nil {
		return nil, err
	}

	// Инициализация сервиса сокращателя ссылок
	err = app.initShortener()
	if err != nil {
		return nil, err
	}

	// Инициализация HTTP сервера
	err = app.initHTTPServer()
	if err != nil {
		return nil, err
	}

	return app, nil
}

// initConfig инициирует конфигурацию приложения.
func (a *App) initConfig() error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	a.cfg = cfg

	return nil
}

// initLogger инициализирует логер.
func (a *App) initLogger() error {
	err := logger.Initialize()
	if err != nil {
		return err
	}

	return nil
}

// initShortener инициализирует сервис сокращателя ссылок, включая хранилище.
func (a *App) initShortener() error {
	repo, err := repository.NewRepository(a.cfg)
	if err != nil {
		return err
	}

	shortener, err := services.NewShortener(repo, a.cfg)
	if err != nil {
		return err
	}

	a.shortener = shortener

	return nil
}

// initHTTPServer инициализирует HTTP сервер.
func (a *App) initHTTPServer() error {
	server, err := url.NewServer(a.shortener, a.cfg)
	if err != nil {
		return err
	}
	a.server = server

	return nil
}

// Run вызывает запуск приложения.
func (a *App) Run() error {
	return a.runShortenerApp()
}

// runShortenerApp запускает приложение.
func (a *App) runShortenerApp() error {
	// Создаем каналы для Graceful Shutdown
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	// Сигнал прерывания
	signal.Notify(quit, os.Interrupt)

	// Graceful Shutdown выполняем в горутине
	go func() {
		<-quit
		slog.Info("shutting down HTTP server")

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// Выключаем HTTP сервер
		if err := a.server.Shutdown(ctx); err != nil {
			slog.Error("HTTP server shutdown error", slog.String("error", err.Error()))
		}

		// Выключаем сервис сокращателя, включая закрытие хранилища
		if err := a.shortener.Close(); err != nil {
			slog.Error("failed to close shortener service", slog.String("error", err.Error()))
		}

		close(done)
	}()

	slog.Info("starting HTTP server", "addr", a.server.Addr)

	// Запускаем HTTP сервер
	var err error
	if a.cfg.EnableHTTPS {
		slog.Info("creating certificate")
		err = cert.CreateCertificate(cert.CertPEM, cert.PrivateKeyPEM)
		if err != nil {
			slog.Error("certificate creation", slog.String("error", err.Error()))
			return err
		}
		err = a.server.ListenAndServeTLS(cert.CertPEM, cert.PrivateKeyPEM)
	} else {
		err = a.server.ListenAndServe()
	}
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("HTTP server error", slog.String("error", err.Error()))
		return err
	}

	<-done
	slog.Info("HTTP server stopped")
	return nil
}
