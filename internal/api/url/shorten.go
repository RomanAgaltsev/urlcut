package url

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
)

// Shorten выполняет обработку запроса на сокращение URL, который передается в текстовом формате.
func (h *Handlers) Shorten(w http.ResponseWriter, r *http.Request) {
	// Читаем оригинальный URL из тела запроса
	longURL, _ := io.ReadAll(r.Body)
	defer func() { _ = r.Body.Close() }()

	// Если оригинального URL нет, считаем, что запрос плохой
	if len(longURL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Выполняем сокращение полученного оригинального URL
	url, err := h.shortener.Shorten(string(longURL))
	if err != nil && !errors.Is(err, repository.ErrConflict) {
		slog.Info(
			"failed to short URL",
			"error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	// Из полученной структуры формируем сокращенный URL
	shortURL := url.Short()

	w.Header().Set("Content-Type", ContentTypeText)
	w.Header().Set("Content-Length", strconv.Itoa(len(shortURL)))

	// Проверяем на наличие конфликта данных - повторная отправка оригинального URL
	// Если конфликт есть, то все равно возвращаем сокращенный URL - возвращается уже имеющийся в БД
	// разница только в статусе ответа
	if errors.Is(err, repository.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	// Пишем сокращенный URL в тело ответа
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		slog.Info(
			"failed to write shorten URL to response",
			"error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
	}
}

// ShortenAPI выполняет обработку запроса на сокращение URL, который передается в формате JSON.
func (h *Handlers) ShortenAPI(w http.ResponseWriter, r *http.Request) {
	// Читать тело запроса будем при помощи JSON декодера
	dec := json.NewDecoder(r.Body)
	defer func() { _ = r.Body.Close() }()

	// Читаем тело запроса
	var req model.Request
	if err := dec.Decode(&req); err != nil {
		slog.Info(
			"failed to unmarshal long URL",
			"error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	// Проверяем, передан ли оригинальный URL
	longURL := req.URL
	if len(longURL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Выполняем сокращение URL
	url, errShort := h.shortener.Shorten(longURL)
	if errShort != nil && !errors.Is(errShort, repository.ErrConflict) {
		slog.Info(
			"failed to short URL",
			"error", errShort.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	// Сокращенный URL преобразуем в JSON
	res, err := json.Marshal(model.Response{Result: url.Short()})
	if err != nil {
		slog.Info(
			"failed to marshal shorten URL",
			"error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	// Пишем заголовки
	w.Header().Set("Content-Type", ContentTypeJSON)
	w.Header().Set("Content-Length", strconv.Itoa(len(res)))

	// Проверяем на наличие конфликта данных - повторная отправка оригинального URL
	// Если конфликт есть, то все равно возвращаем сокращенный URL - возвращается уже имеющийся в БД
	// разница только в статусе ответа
	if errors.Is(errShort, repository.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	// Пишем сокращенный URL в тело ответа
	_, err = w.Write(res)
	if err != nil {
		slog.Info(
			"failed to write shorten URL to response",
			"error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
	}
}

// ShortenAPIBatch выполняет обработку запроса на сокращение массива URL (батча), который передается в формате JSON.
func (h *Handlers) ShortenAPIBatch(w http.ResponseWriter, r *http.Request) {
	// Создаем слайс для батча
	batch := make([]model.BatchRequest, 0)

	// Читать тело запроса будем при помощи JSON декодера
	dec := json.NewDecoder(r.Body)
	defer func() { _ = r.Body.Close() }()

	// Читаем открывающую скобку "["
	_, err := dec.Token()
	if err != nil {
		slog.Info(
			"failed to decode batch",
			"error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	// Читаем данные батча
	for dec.More() {
		var batchReq model.BatchRequest
		if err := dec.Decode(&batchReq); err != nil {
			slog.Info(
				"failed to decode batch element",
				"error", err.Error())
			http.Error(w, "please look at logs", http.StatusInternalServerError)
			return
		}
		// Прочитанный элемент батча сохраняем в слайс
		batch = append(batch, batchReq)
	}

	// Если получили пустой батч, то и делать нечего...
	if len(batch) == 0 {
		slog.Info("got empty batch")
		http.Error(w, "please look at logs", http.StatusBadRequest)
		return
	}

	// Сокращаем все URL батча, которые были прочитаны
	batchShortened, err := h.shortener.ShortenBatch(batch)
	if err != nil {
		slog.Info(
			"failed to short URL",
			"error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	// Пишем заголовки
	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)

	// Данные в тело ответа будем записывать при помощи JSON енкодера
	enc := json.NewEncoder(w)

	err = enc.Encode(batchShortened)
	if err != nil {
		slog.Info(
			"failed to encode batch",
			"error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}
}
