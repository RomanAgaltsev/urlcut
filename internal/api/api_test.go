package api

import (
    "io"
    "net/http"
    "net/http/httptest"
    "strconv"
    "strings"
    "testing"

    "github.com/RomanAgaltsev/urlcut/internal/config"
    "github.com/RomanAgaltsev/urlcut/internal/repository"
    "github.com/RomanAgaltsev/urlcut/internal/service"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// assertEqualBadRequest - проверяет корректность полученного ответа со статусом 400
func assertEqualBadRequest(t *testing.T, actual *http.Response) {
    // Структура ожидаемого ответа
    expected := struct {
        contentType   string // Заголовок "Content-Type"
        contentLength int    // Заголовок "Content-Length"
        body          string // Тело ответа
    }{"text/plain", 11, "Bad request"}

    // Считываем тело ответа
    defer actual.Body.Close()
    resBody, err := io.ReadAll(actual.Body)
    // Проверяем отсутствие ошибок
    require.NoError(t, err)

    // Рассчитываем длину ответа для проверки
    contentLength, _ := strconv.Atoi(actual.Header.Get("Content-Length"))

    // Проверяем заголовок "Content-Type"
    assert.Equal(t, expected.contentType, actual.Header.Get("Content-Type"))
    // Проверяем длину ответа
    assert.Equal(t, expected.contentLength, contentLength)
    // Проверяем тело ответа
    assert.Equal(t, expected.body, string(resBody))
}

func TestShortenHandler(t *testing.T) {
    // Структура отправляемого запроса
    type request struct {
        method      string // Метод запроса
        contentType string // Заголовок "Content-Type"
        body        string // Тело запроса
    }
    // Структура получаемого ответа
    type response struct {
        statusCode    int    // Код статус ответа
        contentType   string // Заголовок "Content-Type"
        contentLength int    // Заголовок "Content-Length"
        body          string // Тело ответа
    }

    // Тесты
    tests := []struct {
        name string
        req  request
        res  response
    }{
        {"[POST] [text/plain] [https://practicum.yandex.ru/]",
            request{http.MethodPost, "text/plain", "https://practicum.yandex.ru/"},
            response{http.StatusCreated, "text/plain", 30, ""},
        },
        {"[POST] [text/plain] [https://translate.yandex.ru/]",
            request{http.MethodPost, "text/plain", "https://translate.yandex.ru/"},
            response{http.StatusCreated, "text/plain", 30, ""},
        },
        {"[POST] [text/plain] ['']",
            request{http.MethodPost, "text/plain", ""},
            response{http.StatusBadRequest, "", 0, ""},
        },
        {"[GET] [text/plain] ['']",
            request{http.MethodGet, "text/plain", ""},
            response{http.StatusBadRequest, "", 0, ""},
        },
        {"[PUT] [text/plain] ['']",
            request{http.MethodGet, "text/plain", ""},
            response{http.StatusBadRequest, "", 0, ""},
        },
    }

    // Создаем структуру конфигурации
    cfg := &config.Config{
        ServerPort: "localhost:8080",
        BaseURL:    "http://localhost:8080",
        IDlength:   8,
    }

    // Чтобы добраться до хендлеров, создаем репо и сервис
    mapRepository := repository.NewMap()
    shortenerService := service.NewShortener(mapRepository, cfg)
    handler := NewHandler(shortenerService, cfg)

    // Запуск тестов
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {

            // Создаем новый запрос
            req := httptest.NewRequest(test.req.method, "/", strings.NewReader(test.req.body))
            // Устанавливаем заголовок "Content-Type" запроса
            req.Header.Set("Content-Type", test.req.contentType)

            // Создаем новый ResponseRecorder
            w := httptest.NewRecorder()
            // Вызваем хендлер запроса на сокращение URL
            handler.ShortenURL(w, req)

            // Получаем результат-ответ
            res := w.Result()

            // Проверяем статус ответа
            assert.Equal(t, test.res.statusCode, res.StatusCode)
            // Если код статуса ответа = 400, обрабатываем отдельно
            if res.StatusCode == http.StatusBadRequest {
                assertEqualBadRequest(t, res)
                return
            }

            // Получаем значение заголовка "Content-Length"
            contentLength, _ := strconv.Atoi(res.Header.Get("Content-Length"))

            // Проверяем заголовок "Content-Type"
            assert.Equal(t, test.res.contentType, res.Header.Get("Content-Type"))
            // Проверяем заголовок "Content-Length"
            assert.Equal(t, test.res.contentLength, contentLength)
        })
    }
}

func TestExpandHandler(t *testing.T) {
    // Структура отправляемого запроса
    type request struct {
        method      string // Метод запроса
        url         string // URL для сокращения
        contentType string // Заголовок "Content-Type"
    }
    // Структура получаемого ответа
    type response struct {
        statusCode int    // Код статус ответа
        header     string // Имя заголовка для проверки (Location)
        url        string // URL для проверки
    }

    // Тесты
    tests := []struct {
        name string
        req  request
        res  response
    }{
        {"[GET] [https://practicum.yandex.ru/] [text/plain]",
            request{http.MethodGet, "https://practicum.yandex.ru/", "text/plain"},
            response{http.StatusTemporaryRedirect, "Location", "https://practicum.yandex.ru/"},
        },
        {"[GET] [https://translate.yandex.ru/] [text/plain]",
            request{http.MethodGet, "https://translate.yandex.ru/", "text/plain"},
            response{http.StatusTemporaryRedirect, "Location", "https://translate.yandex.ru/"},
        },
        {"[POST] [https://translate.yandex.ru/] [text/plain]",
            request{http.MethodPost, "https://translate.yandex.ru/", "text/plain"},
            response{http.StatusBadRequest, "", ""},
        },
    }

    // Создаем структуру конфигурации
    cfg := &config.Config{
        ServerPort: "localhost:8080",
        BaseURL:    "http://localhost:8080",
        IDlength:   8,
    }

    // Чтобы добраться до хендлеров, создаем репо и сервис
    mapRepository := repository.NewMap()
    shortenerService := service.NewShortener(mapRepository, cfg)
    handler := NewHandler(shortenerService, cfg)

    // Запуск тестов
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            // Сначала создаем POST запрос на сокращение URL и получения идентификатора сокращенного URL
            reqPost := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.req.url))
            // Устанавливаем заголовок "Content-Type" запроса
            reqPost.Header.Set("Content-Type", test.req.contentType)

            // Создаем новый ResponseRecorder
            wPost := httptest.NewRecorder()
            // Вызваем хендлер запроса на сокращение URL
            handler.ShortenURL(wPost, reqPost)

            // Получаем результат-ответ
            resPost := wPost.Result()

            // Откладываем закрытие тела ответа
            defer resPost.Body.Close()
            // Получаем тело ответа
            resPostBody, err := io.ReadAll(resPost.Body)
            // Проверяем отсутствие ошибок при чтении тела
            require.NoError(t, err)

            // Получаем сокращенный URL из тела ответа
            shortenedURL := string(resPostBody)
            // Получаем идентификатор сокращенного URL
            urlID := strings.TrimPrefix(shortenedURL, cfg.BaseURL+"/")

            // Создаем новый запрос на получение оригинального URL по идентификатору сокращенного
            req := httptest.NewRequest(test.req.method, "/"+urlID, nil)
            // Устанавливаем заголовок "Content-Type" запроса
            req.Header.Set("Content-Type", test.req.contentType)

            // Создаем новый ResponseRecorder
            w := httptest.NewRecorder()
            // Вызываем хендлер запроса на получение оригинального URL
            handler.ExpandURL(w, req)

            // Получаем результат-ответ
            res := w.Result()

            // Проверяем статус ответа
            assert.Equal(t, test.res.statusCode, res.StatusCode)
            // Если код статуса ответа = 400, обрабатываем отдельно
            if res.StatusCode == http.StatusBadRequest {
                assertEqualBadRequest(t, res)
                return
            }

            // Проверяем, содержит ли ответ заголовок - Location
            if assert.Contains(t, res.Header, test.res.header) {
                // Если заголовок есть, проверяем его содержимое - оригинальные URL
                assert.Equal(t, test.res.url, res.Header.Get(test.res.header))
            }
        })
    }
}
