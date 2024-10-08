package url

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

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
