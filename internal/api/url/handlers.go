package url

import (
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	servicesurl "github.com/RomanAgaltsev/urlcut/internal/services/url"
)

// Handlers - структура для хендлеров сервиса
type Handlers struct {
	service servicesurl.Service // Сервис сокращения URL
}

// NewHandlers - конструктор структуры обработчиков
func NewHandlers(service servicesurl.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

// ShortenURL - хендлер сокращения URL
func (h *Handlers) ShortenURL(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса - принимаем только POST
	if r.Method != http.MethodPost {
		// Это не POST запрос, возвращаем статус 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Получаем URL из тела запроса
	url, _ := io.ReadAll(r.Body)
	// Откладываем закрытие тела запроса
	defer func() { _ = r.Body.Close() }()
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

// ExpandURL - хендлер возврата оригинального URL
func (h *Handlers) ExpandURL(w http.ResponseWriter, r *http.Request) {
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

type (
	// responseData - структура для хранения данных ответа
	responseData struct {
		status int // Статус ответа
		size   int // Размер содержимого ответа
	}

	// loggingResponseWriter - структура-обертка для ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter               // Встраиваем ResponseWriter
		responseData        *responseData // Данные ответа
	}
)

// Write - реализует Write с подсчетом размера данных
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader - реализует WriteHeader с регистрацией кода статуса ответа
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// WithLogging - middleware с подключением логирования для хендлеров
func WithLogging(h http.HandlerFunc) http.HandlerFunc {
	// Возвращать будем новую функцию-хендлер
	logFn := func(w http.ResponseWriter, r *http.Request) {
		// Фиксируем время старта обработки
		start := time.Now()

		// Создаем структуру для хранения данных ответа
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		// Подменяем полученный ResponseWriter собственным, с логированием
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		// Обрабатываем запрос
		h.ServeHTTP(&lw, r)

		// Фиксируем длительность обработки запроса
		duration := time.Since(start)

		// Пишем в лог данные запроса и ответа
		slog.Info(
			"got request",
			slog.String("uri", r.RequestURI),        // URI запроса
			slog.String("method", r.Method),         // Метод запроса
			slog.Int("status", responseData.status), // Статус ответа
			slog.Duration("duration", duration),     // Длительность обработки запроса
			slog.Int("size", responseData.size),     // Размер содержимого ответа
		)
	}
	// Возвращаем новую функцию-хендлер
	return logFn
}
