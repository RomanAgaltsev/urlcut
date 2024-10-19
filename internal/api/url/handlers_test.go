package url

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	servicesurl "github.com/RomanAgaltsev/urlcut/internal/services/url"

	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/logger"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenHandler(t *testing.T) {
	tests := []struct {
		name      string
		reqMethod string
		reqURL    string
		resStatus int
	}{
		{"[POST] [https://practicum.yandex.ru/]",
			http.MethodPost, "https://practicum.yandex.ru/", http.StatusCreated},
		{"[POST] [https://translate.yandex.ru/]",
			http.MethodPost, "https://practicum.yandex.ru/", http.StatusCreated},
		{"[POST] ['']",
			http.MethodPost, "", http.StatusBadRequest},
		{"[GET] ['']",
			http.MethodGet, "", http.StatusBadRequest},
		{"[PUT] ['']",
			http.MethodPut, "", http.StatusBadRequest},
	}

	// Создаем структуру конфигурации
	cfg := &config.Config{
		ServerPort: "localhost:8080",
		BaseURL:    "http://localhost:8080",
		IDlength:   8,
	}

	_ = logger.Initialize()

	// Чтобы добраться до хендлеров, создаем репо и сервис
	repo := repository.New()
	service := servicesurl.NewShortener(repo, cfg.BaseURL, cfg.IDlength)
	handlers := NewHandlers(service)

	// Запуск тестов
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Создаем новый запрос
			req := httptest.NewRequest(test.reqMethod, "/", strings.NewReader(test.reqURL))
			// Устанавливаем заголовок "Content-Type" запроса
			req.Header.Set("Content-Type", "text/plain")

			// Создаем новый ResponseRecorder
			w := httptest.NewRecorder()
			// Вызваем хендлер запроса на сокращение URL
			handlers.ShortenURL(w, req)

			// Получаем результат-ответ
			res := w.Result()

			// Проверяем статус ответа
			assert.Equal(t, test.resStatus, res.StatusCode)
			if test.resStatus == http.StatusBadRequest {
				return
			}

			// Проверяем заголовок "Content-Type"
			assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))

			// Откладываем закрытие тела ответа
			defer func() { _ = res.Body.Close() }()
			// Получаем тело ответа
			resBody, err := io.ReadAll(res.Body)
			// Проверяем отсутствие ошибок при чтении тела
			require.NoError(t, err)

			// Получаем сокращенный URL из тела ответа
			shortenedURL := string(resBody)
			// Проверяем префикс
			assert.Equal(t, strings.HasPrefix(shortenedURL, cfg.BaseURL), true)
		})
	}
}

func TestExpandHandler(t *testing.T) {
	tests := []struct {
		name      string
		reqMethod string
		reqURL    string
		resStatus int
	}{
		{"[GET] [https://practicum.yandex.ru/]",
			http.MethodGet, "https://practicum.yandex.ru/", http.StatusTemporaryRedirect},
		{"[GET] [https://translate.yandex.ru/]",
			http.MethodGet, "https://translate.yandex.ru/", http.StatusTemporaryRedirect},
		{"[POST] [https://translate.yandex.ru/]",
			http.MethodPost, "https://translate.yandex.ru/", http.StatusBadRequest},
	}

	// Создаем структуру конфигурации
	cfg := &config.Config{
		ServerPort: "localhost:8080",
		BaseURL:    "http://localhost:8080",
		IDlength:   8,
	}

	// Чтобы добраться до хендлеров, создаем репо и сервис
	repo := repository.New()
	service := servicesurl.NewShortener(repo, cfg.BaseURL, cfg.IDlength)
	handlers := NewHandlers(service)

	// Запуск тестов
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Сначала создаем POST запрос на сокращение URL и получения идентификатора сокращенного URL
			reqPost := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.reqURL))
			// Устанавливаем заголовок "Content-Type" запроса
			reqPost.Header.Set("Content-Type", "text/plain")

			// Создаем новый ResponseRecorder
			wPost := httptest.NewRecorder()
			// Вызваем хендлер запроса на сокращение URL
			handlers.ShortenURL(wPost, reqPost)

			// Получаем результат-ответ
			resPost := wPost.Result()

			// Откладываем закрытие тела ответа
			defer func() { _ = resPost.Body.Close() }()
			// Получаем тело ответа
			resPostBody, err := io.ReadAll(resPost.Body)
			// Проверяем отсутствие ошибок при чтении тела
			require.NoError(t, err)

			// Получаем сокращенный URL из тела ответа
			shortenedURL := string(resPostBody)
			// Получаем идентификатор сокращенного URL
			urlID := strings.TrimPrefix(shortenedURL, cfg.BaseURL+"/")

			// Создаем новый запрос на получение оригинального URL по идентификатору сокращенного
			req := httptest.NewRequest(test.reqMethod, "/"+urlID, nil)
			// Устанавливаем заголовок "Content-Type" запроса
			req.Header.Set("Content-Type", "text/plain")

			// Создаем новый ResponseRecorder
			w := httptest.NewRecorder()
			// Вызываем хендлер запроса на получение оригинального URL
			handlers.ExpandURL(w, req)

			// Получаем результат-ответ
			res := w.Result()
			defer func() { _ = res.Body.Close() }()

			// Проверяем статус ответа
			assert.Equal(t, test.resStatus, res.StatusCode)
			// Если код статуса ответа = 400, обрабатываем отдельно
			if res.StatusCode == http.StatusBadRequest {
				return
			}

			// Проверяем, содержит ли ответ заголовок - Location
			if assert.Contains(t, res.Header, "Location") {
				// Если заголовок есть, проверяем его содержимое - оригинальные URL
				assert.Equal(t, test.reqURL, res.Header.Get("Location"))
			}
		})
	}
}
