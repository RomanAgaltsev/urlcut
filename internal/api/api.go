package api

import (
    "github.com/RomanAgaltsev/urlcut/internal/config"
    "github.com/RomanAgaltsev/urlcut/internal/service"
    "github.com/go-chi/chi/v5"
    "io"
    "log"
    "net/http"
    "strconv"
    "strings"
)

// Handler - структура, которая реализует работу с запросами и ответами к серверу
type Handler struct {
    serverPort string          // Адрес сервера и его порт
    service    service.Service // Сервис сокращения URL
    router     chi.Router      // Роутер
}

// NewHandler - конструктор для Handler
func NewHandler(service service.Service, cfg *config.Config) *Handler {
    // Создаем структуру и заполняем из конфигурации
    h := &Handler{
        serverPort: cfg.ServerPort,
        service:    service,
        router:     chi.NewRouter(),
    }
    // Устанавливаем хендлеры
    h.SetRoutes()
    return h
}

func (h *Handler) SetRoutes() {
    // Добавляем хендлеры
    h.router.Post("/", h.ShortenURL)   // Запрос на сокращение URL - POST
    h.router.Get("/{id}", h.ExpandURL) // Запрос на возврат исходного URL - GET
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
    // Проверяем метод запроса - принимаем только POST
    if r.Method != http.MethodPost {
        // Это не POST запрос, возвращаем статус 400
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    // Получаем URL из тела запроса
    url, _ := io.ReadAll(r.Body)
    // Откладываем закрытие тела запроса
    defer r.Body.Close()
    // Проверяем, передали ли URL
    if len(url) == 0 {
        // URL не передали, возвращаем статус 400
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    // Формируем сокращенный URL
    shortenedURL, err := h.service.ShortenURL(string(url))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // Пишем заголовки в ответ
    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Content-Length", strconv.Itoa(len(shortenedURL)))
    // Пишем статус 201 в ответ
    w.WriteHeader(http.StatusCreated)
    // Пишем сокращенный URL в ответ
    _, err = w.Write([]byte(shortenedURL))
    if err != nil {
        log.Printf("writing of shortened URL failed: %v", err)
    }
}

func (h *Handler) ExpandURL(w http.ResponseWriter, r *http.Request) {
    // Проверяем метод запроса - принимаем только GET
    if r.Method != http.MethodGet {
        // Это не GET запрос, возвращаем статус 400
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    // Удалем префикс из полученного идентификатора
    urlID := strings.TrimPrefix(r.URL.Path, "/")
    // Получаем оригинальный URL
    expandedURL, err := h.service.ExpandURL(urlID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
    }
    // Пишем заголовки в ответ
    w.Header().Set("Location", expandedURL)
    // Пишем статус 307 в ответ
    w.WriteHeader(http.StatusTemporaryRedirect)
}

// Run - Запускает HTTP-сервер
func (h *Handler) Run() {
    log.Fatal(http.ListenAndServe(h.serverPort, h.router))
}
