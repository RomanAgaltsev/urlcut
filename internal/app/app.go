package app

import (
	"fmt"
	"log"
	"net/http"

	apiurl "github.com/RomanAgaltsev/urlcut/internal/api/url"
	servicesurl "github.com/RomanAgaltsev/urlcut/internal/services/url"

	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
	"github.com/go-chi/chi/v5"
)

// App - структура приложения
type App struct {
	repo    repository.Repository // Репозиторий URL
	service servicesurl.Service   // Сервис сокращения URL
	server  *http.Server          // HTTP-сервер
}

// New - возвращает новый экземпляр приложения
func New() (*App, error) {
	// Получаем конфигурацию
	cfg, err := config.Get()
	// Проверяем наличие ошибки
	if err != nil {
		// Есть ошибка, возвращаем nil и ошибку
		return nil, fmt.Errorf("getting config failed: %v", err)
	}
	// Создаем новое приложение
	app := &App{}
	// Получаем репозиторий
	err = app.getRepository()
	// Проверяем наличие ошибки
	if err != nil {
		// Есть ошибка, возвращаем nil и ошибку
		return nil, err
	}
	// Получаем сервис
	err = app.getService(cfg.BaseURL, cfg.IDlength)
	// Проверяем наличие ошибки
	if err != nil {
		// Есть ошибка, возвращаем nil и ошибку
		return nil, err
	}
	// Получаем HTTP-сервер
	err = app.getHTTPServer(cfg.ServerPort)
	// Проверяем наличие ошибки
	if err != nil {
		// Есть ошибка, возвращаем nil и ошибку
		return nil, err
	}
	// Ошибок не было, возвращаем приложение
	return app, nil
}

// getRepository - устанавливает репозиторий в приложении
func (a *App) getRepository() error {
	a.repo = repository.New()
	return nil
}

// getService - устанавливает сервис сокращения URL в приложении
func (a *App) getService(baseURL string, idLength int) error {
	a.service = servicesurl.NewShortener(a.repo, baseURL, idLength)
	return nil
}

// getHTTPServer - устанавливает HTTP-сервер в приложении
func (a *App) getHTTPServer(serverPort string) error {
	// Получаем обработчики
	handlers := apiurl.NewHandlers(a.service)
	// Создаем новый роутер
	router := chi.NewRouter()
	// Добавляем хендлеры
	router.Post("/", handlers.ShortenURL)   // Запрос на сокращение URL - POST
	router.Get("/{id}", handlers.ExpandURL) // Запрос на возврат исходного URL - GET
	// Создаем новый HTTP-сервер
	a.server = &http.Server{
		Addr:    serverPort,
		Handler: router,
	}
	return nil
}

// Run - запускает приложение
func (a *App) Run() {
	a.runShortenerApp()
}

// runShortenerApp - запускает HTTP-сервер
func (a *App) runShortenerApp() {
	if err := a.server.ListenAndServe(); err != nil {
		log.Fatalf("running HTTP server failed: %s", err.Error())
	}
}
