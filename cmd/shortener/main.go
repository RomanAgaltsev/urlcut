package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/RomanAgaltsev/urlcut/internal/config"

	"github.com/go-chi/chi/v5"
)

var urls map[string]string // Переменная-хранилище, в которой хранится полученный URL

func init() {
	// Инициируем хранилище
	urls = make(map[string]string)
}

// urlID - возвращает идентификатор сокращенного URL
func urlID() string {
	const (
		// Символы, которые могут входить в идентификатор
		letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	)
	// Инициируем слайс байт с длиной, равной длине идентификатора
	//b := make([]byte, length)
	b := make([]byte, config.Config.IDlength)
	// Заполняем слайс произвольными символами из доступных
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	// Возвращаем получившуюся строку - идентификатор
	return string(b)
}

// shortenURL - возвращает сокращенный URL для переданного URL
func shortenURL(url string) string {
	// Получаем новый идентификатор
	id := urlID()
	// Сохраняем пару идентификатор-URL
	urls[id] = url
	// Возвращаем сокращенный URL, включая адрес сервера
	return config.Config.BasicAddr + "/" + id
}

// expandURL - Возвращает оригинальный URL по переданному идентификатору
func expandURL(urlID string) string {
	// Возвращаем, если есть что возвращать
	return urls[urlID]
}

// shortenHandler - хендлер для обработки запроса на сокращение URL
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса - принимаем только POST
	if r.Method != http.MethodPost {
		// Это не POST запрос, возвращаем статус 400
		badRequestHandler(w)
		return
	}
	// Получаем URL из тела запроса
	url, _ := io.ReadAll(r.Body)
	// Откладываем закрытие тела запроса
	defer r.Body.Close()
	// Проверяем, передали ли URL
	if len(url) == 0 {
		// URL не передали, возвращаем статус 400
		badRequestHandler(w)
		return
	}
	// Формируем сокращенный URL
	shortenedURL := shortenURL(string(url))
	// Пишем заголовки в ответ
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(shortenedURL)))
	// Пишем статус 201 в ответ
	w.WriteHeader(http.StatusCreated)
	// Пишем сокращенный URL в ответ
	w.Write([]byte(shortenedURL))
}

// expandHandler - хендлер для обработки запроса на возврат сходного URL
func expandHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса - принимаем только GET
	if r.Method != http.MethodGet {
		// Это не GET запрос, возвращаем статус 400
		badRequestHandler(w)
		return
	}
	// Удалем префикс из полученного идентификатора
	urlID := strings.TrimPrefix(r.URL.Path, "/")
	// Получаем оригинальный URL
	expandedURL := expandURL(urlID)
	// Пишем заголовки в ответ
	w.Header().Set("Location", string(expandedURL))
	// Пишем статус 307 в ответ
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// badRequestHandler - универсальный обработчик для возврата статуса 400
func badRequestHandler(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len("Bad request")))
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad request"))
}

func main() {
	config.ParseFlags()

	// Создаем новый роутер
	r := chi.NewRouter()

	// Добавляем хендлеры
	r.Post("/", shortenHandler)   // Запрос на сокращение URL - POST
	r.Get("/{id}", expandHandler) // Запрос на возврат исходного URL - GET

	// Запускаем сервер
	log.Fatal(http.ListenAndServe(config.Config.ServerAddr, r))
}
