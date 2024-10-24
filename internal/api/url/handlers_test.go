package url

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/urlcut/internal/api/middleware"
	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/logger"
	"github.com/RomanAgaltsev/urlcut/internal/model"
	repositoryurl "github.com/RomanAgaltsev/urlcut/internal/repository/url"
	serviceurl "github.com/RomanAgaltsev/urlcut/internal/service/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
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
			http.MethodGet, "", http.StatusMethodNotAllowed},
		{"[PUT] ['']",
			http.MethodPut, "", http.StatusMethodNotAllowed},
	}

	cfg := &config.Config{
		ServerPort: "localhost:8080",
		BaseURL:    "http://localhost:8080",
		IDlength:   8,
	}

	_ = logger.Initialize()

	repo := repositoryurl.New()
	service := serviceurl.New(repo, cfg.BaseURL, cfg.IDlength)
	handlers := New(service)

	router := chi.NewRouter()
	router.Use(middleware.WithLogging)
	router.Use(middleware.WithGzip)
	router.Post("/", handlers.Shorten)

	httpSrv := httptest.NewServer(router)
	defer httpSrv.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = test.reqMethod
			req.URL = httpSrv.URL

			res, err := req.
				SetHeader("Content-Type", "text/plain").
				SetBody(test.reqURL).
				Send()
			assert.NoError(t, err)

			assert.Equal(t, test.resStatus, res.StatusCode())
			if test.resStatus == http.StatusBadRequest || test.resStatus == http.StatusMethodNotAllowed {
				return
			}

			assert.Equal(t, "text/plain", res.Header().Get("Content-Type"))

			shortenedURL := string(res.Body())
			assert.Equal(t, strings.HasPrefix(shortenedURL, cfg.BaseURL), true)
		})
	}
}

func TestShortenAPIHandler(t *testing.T) {
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
			http.MethodGet, "", http.StatusMethodNotAllowed},
		{"[PUT] ['']",
			http.MethodPut, "", http.StatusMethodNotAllowed},
	}

	cfg := &config.Config{
		ServerPort: "localhost:8080",
		BaseURL:    "http://localhost:8080",
		IDlength:   8,
	}

	_ = logger.Initialize()

	repo := repositoryurl.New()
	service := serviceurl.New(repo, cfg.BaseURL, cfg.IDlength)
	handlers := New(service)

	router := chi.NewRouter()
	router.Use(middleware.WithLogging)
	router.Use(middleware.WithGzip)
	router.Post("/api/shorten", handlers.ShortenAPI)

	httpSrv := httptest.NewServer(router)
	defer httpSrv.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = test.reqMethod
			req.URL = httpSrv.URL + "/api/shorten"

			request := model.Request{
				URL: test.reqURL,
			}
			reqBytes, _ := json.Marshal(request)

			res, err := req.
				SetHeader("Content-Type", "application/json").
				SetHeader("Accept-Encoding", "gzip").
				//			SetHeader("Content-Encoding", "gzip").
				SetBody(bytes.NewReader(reqBytes)).
				Send()
			assert.NoError(t, err)

			assert.Equal(t, test.resStatus, res.StatusCode())
			if test.resStatus == http.StatusBadRequest || test.resStatus == http.StatusMethodNotAllowed {
				return
			}
			assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

			dec := json.NewDecoder(bytes.NewReader(res.Body()))
			var response model.Response
			err = dec.Decode(&response)
			require.NoError(t, err)

			shortenedURL := response.Result
			assert.Equal(t, strings.HasPrefix(shortenedURL, cfg.BaseURL), true)
		})
	}

	t.Run("[POST] [nil body]", func(t *testing.T) {
		req := resty.New().R()
		req.Method = http.MethodPost
		req.URL = httpSrv.URL + "/api/shorten"

		res, err := req.
			SetHeader("Content-Type", "application/json").
			Send()
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode())
	})

	t.Run("[POST] [Bad body]", func(t *testing.T) {
		request := struct {
			Dummy string
		}{
			Dummy: "Hi there!",
		}
		reqBytes, _ := json.Marshal(request)

		req := resty.New().R()
		req.Method = http.MethodPost
		req.URL = httpSrv.URL + "/api/shorten"

		res, err := req.
			SetHeader("Content-Type", "application/json").
			SetBody(bytes.NewReader(reqBytes)).
			Send()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})
}

func TestExpandHandler(t *testing.T) {
	tests := []struct {
		name      string
		reqMethod string
		reqURL    string
		resStatus int
	}{
		{"[GET] [https://practicum.yandex.ru/]",
			http.MethodGet, "https://practicum.yandex.ru/", http.StatusOK},
		{"[GET] [https://translate.yandex.ru/]",
			http.MethodGet, "https://translate.yandex.ru/", http.StatusOK},
		{"[POST] [https://translate.yandex.ru/]",
			http.MethodPost, "https://translate.yandex.ru/", http.StatusMethodNotAllowed},
	}

	cfg := &config.Config{
		ServerPort: "localhost:8080",
		BaseURL:    "http://localhost:8080",
		IDlength:   8,
	}

	_ = logger.Initialize()

	repo := repositoryurl.New()
	service := serviceurl.New(repo, cfg.BaseURL, cfg.IDlength)
	handlers := New(service)

	router := chi.NewRouter()
	router.Use(middleware.WithLogging)
	router.Use(middleware.WithGzip)
	router.Post("/", handlers.Shorten)
	router.Get("/{id}", handlers.Expand)

	httpSrv := httptest.NewServer(router)
	defer httpSrv.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reqPost := resty.New().R()
			reqPost.Method = http.MethodPost
			reqPost.URL = httpSrv.URL

			resPost, err := reqPost.
				SetHeader("Content-Type", "text/plain").
				SetBody(test.reqURL).
				Send()
			assert.NoError(t, err)

			shortenedURL := string(resPost.Body())
			urlID := strings.TrimPrefix(shortenedURL, cfg.BaseURL+"/")

			req := resty.New().R()
			req.Method = test.reqMethod
			req.URL = httpSrv.URL + "/" + urlID

			res, err := req.
				SetHeader("Content-Type", "text/plain").
				Send()
			assert.NoError(t, err)

			assert.Equal(t, test.resStatus, res.StatusCode())
			if res.StatusCode() == http.StatusBadRequest || res.StatusCode() == http.StatusMethodNotAllowed {
				return
			}

//			if assert.Contains(t, res.Header(), "Location") {
//				assert.Equal(t, test.reqURL, res.Header().Get("Location"))
//			}
		})
	}
}
