package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	serverAddr   string = ":8080"                          // Адрес сервера, пока, храним в константе
	shortenedURL string = "http://localhost:8080/EwHXdJfB" // Всегда возвращаем один и тот же URL - это временно
)

var url []byte // Переменная-хранилище, в которой хранится полученный URL

// shortenHandler - хендлер для обработки запроса на сокращение URL
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса - принимаем только POST
	if r.Method != http.MethodPost {
		// Это не POST запрос, возвращаем статус 400
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}
	// Получаем URL из тела запроса
	url, _ = io.ReadAll(r.Body)
	// Откладываем закрытие тела запроса
	defer r.Body.Close()
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
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}
	// Пишем заголовки в ответ
	w.Header().Set("Location", string(url))
	// Пишем статус 307 в ответ
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	// Создаем собственный обработчик/мультиплексор
	mux := http.NewServeMux()
	// Добавляем хендлеры
	mux.HandleFunc("/", shortenHandler)    // Запрос на сокращение URL
	mux.HandleFunc("/{id}", expandHandler) // Запрос на возврат исходного URL
	// Запускаем сервер
	log.Fatal(http.ListenAndServe(serverAddr, mux))
}
