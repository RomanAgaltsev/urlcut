// Пакет url реализует http-сервер и хендлеры получаемых сервером запросов.
package url

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/database"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/pkg/auth"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
)

// Служебные константы для заголовков.
const (
	// ContentTypeJSON используется для установки значений заголовков http-ответов.
	ContentTypeJSON = "application/json"

	// ContentTypeText используется для установки значений заголовков http-ответов.
	ContentTypeText = "text/plain; charset=utf-8"
)

// ErrNoUserID ошибка отсутствия идентификатора пользователя в полученном запросе (в куке или заголовке).
var ErrNoUserID = fmt.Errorf("no user ID provided")

// Handlers владеет хендлерами получаемых http-запросов.
type Handlers struct {
	shortener interfaces.Service // сервис сокращателя ссылок
	cfg       *config.Config     // конфигурация приложения
}

// NewHandlers - функция-конструктор Handlers.
func NewHandlers(shortener interfaces.Service, cfg *config.Config) *Handlers {
	return &Handlers{
		shortener: shortener,
		cfg:       cfg,
	}
}

// Shorten выполняет обработку запроса на сокращение URL, который передается в текстовом формате.
// Перед обработкой выполняется извлечение идентификатора пользователя из запроса.
// При отсутствии идентификатора, обработка запроса прекращается и возвращается статус Unauthorized.
func (h *Handlers) Shorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем идентификатор пользователя
	uid, err := getUserUID(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Читаем оригинальный URL из тела запроса
	longURL, _ := io.ReadAll(r.Body)
	defer func() { _ = r.Body.Close() }()

	// Если оригинального URL нет, считаем, что запрос плохой
	if len(longURL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Выполняем сокращение полученного оригинального URL
	url, err := h.shortener.Shorten(ctx, string(longURL), uid)
	if err != nil && !errors.Is(err, repository.ErrConflict) {
		slog.Info("failed to short URL", "error", err.Error())
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
		slog.Info("failed to write shorten URL to response", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
	}
}

// ShortenAPI выполняет обработку запроса на сокращение URL, который передается в формате JSON.
// Перед обработкой выполняется извлечение идентификатора пользователя из запроса.
// При отсутствии идентификатора, обработка запроса прекращается и возвращается статус Unauthorized.
func (h *Handlers) ShortenAPI(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем идентификатор пользователя
	uid, err := getUserUID(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Читать тело запроса будем при помощи JSON декодера
	dec := json.NewDecoder(r.Body)
	defer func() { _ = r.Body.Close() }()

	// Читаем тело запроса
	var req model.URLDTO
	if err = dec.Decode(&req); err != nil {
		slog.Info("failed to unmarshal long URL", "error", err.Error())
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
	url, errShort := h.shortener.Shorten(ctx, longURL, uid)
	if errShort != nil && !errors.Is(errShort, repository.ErrConflict) {
		slog.Info("failed to short URL", "error", errShort.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	// Сокращенный URL преобразуем в JSON
	res, err := json.Marshal(model.ResultDTO{Result: url.Short()})
	if err != nil {
		slog.Info("failed to marshal shorten URL", "error", err.Error())
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
		slog.Info("failed to write shorten URL to response", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
	}
}

// ShortenAPIBatch выполняет обработку запроса на сокращение массива URL (батча), который передается в формате JSON.
// Перед обработкой выполняется извлечение идентификатора пользователя из запроса.
// При отсутствии идентификатора, обработка запроса прекращается и возвращается статус Unauthorized.
func (h *Handlers) ShortenAPIBatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем идентификатор пользователя
	uid, err := getUserUID(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Создаем слайс для батча
	batch := make([]model.IncomingBatchDTO, 0)

	// Читать тело запроса будем при помощи JSON декодера
	dec := json.NewDecoder(r.Body)
	defer func() { _ = r.Body.Close() }()

	// Читаем открывающую скобку "["
	_, err = dec.Token()
	if err != nil {
		slog.Info("failed to decode batch", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	// Читаем данные батча
	for dec.More() {
		var batchReq model.IncomingBatchDTO
		if err = dec.Decode(&batchReq); err != nil {
			slog.Info("failed to decode batch element", "error", err.Error())
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
	batchShortened, err := h.shortener.ShortenBatch(ctx, batch, uid)
	if err != nil && !errors.Is(err, repository.ErrConflict) {
		slog.Info("failed to short URL", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	// Пишем заголовки
	w.Header().Set("Content-Type", ContentTypeJSON)

	// Проверяем на наличие конфликта данных - повторная отправка оригинального URL
	// Если конфликт есть, то все равно возвращаем сокращенный URL - возвращается уже имеющийся в БД
	// разница только в статусе ответа
	if errors.Is(err, repository.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	// Данные в тело ответа будем записывать при помощи JSON енкодера
	enc := json.NewEncoder(w)

	err = enc.Encode(batchShortened)
	if err != nil {
		slog.Info("failed to encode batch", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}
}

// Expand выполняет обработку запроса на получение оригинального URL.
// Перед обработкой выполняется извлечение идентификатора пользователя из запроса.
// При отсутствии идентификатора, обработка запроса прекращается и возвращается статус Unauthorized.
func (h *Handlers) Expand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем идентификатор пользователя
	_, err := getUserUID(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Получаем идентификатор из параметров URL
	urlID := chi.URLParam(r, "id")

	// Получаем URL по идентификатору
	url, err := h.shortener.Expand(ctx, urlID)
	if err != nil {
		slog.Info("failed to expand URL", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusNotFound)
		return
	}

	// Если запросили помеченный на удаление URL, возвращаем статус 410
	if url.Deleted {
		w.WriteHeader(http.StatusGone)
		return
	}

	// По идентификатору ничего не нашли
	if len(url.Long) == 0 {
		http.Error(w, "URL ID was not found in repository", http.StatusNotFound)
		return
	}

	// Пишем заголовки
	w.Header().Set("Location", url.Long)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// Ping выполняет обработку запроса на пинг хранилища.
// Перед обработкой выполняется извлечение идентификатора пользователя из запроса.
// При отсутствии идентификатора, обработка запроса прекращается и возвращается статус Unauthorized.
func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем идентификатор пользователя
	_, err := getUserUID(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	db, err := database.NewConnection(ctx, "pgx", h.cfg.DatabaseDSN)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() { _ = db.Close() }()

	w.WriteHeader(http.StatusOK)
}

// UserUrls выполняет обработку запроса на получение списка всех сохраненных URL пользователя.
// Перед обработкой выполняется извлечение идентификатора пользователя из запроса.
// При отсутствии идентификатора, обработка запроса прекращается и возвращается статус Unauthorized.
func (h *Handlers) UserUrls(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем идентификатор пользователя
	uid, err := getUserUID(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Получаем URL-ы пользователя
	urls, err := h.shortener.UserURLs(ctx, uid)
	if err != nil {
		slog.Info("failed to fetch user URLs", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	// Проверяем, если ли данные для возврата
	if len(urls) == 0 {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	}

	// Пишем заголовки
	w.Header().Set("Content-Type", ContentTypeJSON)

	w.WriteHeader(http.StatusOK)

	// Данные в тело ответа будем записывать при помощи JSON енкодера
	enc := json.NewEncoder(w)

	err = enc.Encode(urls)
	if err != nil {
		slog.Info("failed to encode user URLs", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}
}

// UserUrlsDelete выполняет обработку запроса на удаление всех сохраненных URL пользователя.
// Перед обработкой выполняется извлечение идентификатора пользователя из запроса.
// При отсутствии идентификатора, обработка запроса прекращается и возвращается статус Unauthorized.
func (h *Handlers) UserUrlsDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем идентификатор пользователя
	uid, err := getUserUID(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Читаем массив идентификаторов URL из тела запроса
	urlArray, _ := io.ReadAll(r.Body)
	defer func() { _ = r.Body.Close() }()

	// Парсим JSON
	var shortURLs []string
	err = json.Unmarshal(urlArray, &shortURLs)
	if err != nil {
		slog.Info("failed to unmarshal URL ID array", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	// Проверяем, есть ли идентификаторы
	if len(shortURLs) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Устанавливаем пометки удаления
	err = h.shortener.DeleteUserURLs(ctx, uid, &model.ShortURLsDTO{IDs: shortURLs})
	if err != nil {
		slog.Info("failed to delete user URLs", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// Stats выполняет обработку запроса на получении статистики сервиса.
func (h *Handlers) Stats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем идентификатор пользователя
	_, err := getUserUID(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	stats, err := h.shortener.Stats(ctx)
	if err != nil {
		slog.Info("failed to get stats", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	// Пишем заголовки
	w.Header().Set("Content-Type", ContentTypeJSON)

	w.WriteHeader(http.StatusOK)

	// Данные в тело ответа будем записывать при помощи JSON енкодера
	enc := json.NewEncoder(w)

	err = enc.Encode(stats)
	if err != nil {
		slog.Info("failed to encode stats", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}
}

// getUserUid получает идентификатор пользователя из контекста запроса.
func getUserUID(r *http.Request) (uuid.UUID, error) {
	// Получаем идентификатор-интерфейс пользователя из контекста
	uidInterface := r.Context().Value(auth.UserIDClaimName)
	if uidInterface == nil {
		return uuid.Nil, ErrNoUserID
	}

	// Идентификатор-интерфейс есть, пробуем привести к строке
	uidString, ok := uidInterface.(string)
	if !ok {
		return uuid.Nil, ErrNoUserID
	}

	// Пробуем парсить строку в uuid
	uid, err := uuid.Parse(uidString)
	if err != nil {
		return uuid.Nil, ErrNoUserID
	}

	return uid, nil
}
