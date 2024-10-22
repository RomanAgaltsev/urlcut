package url

import (
	"bytes"
	"encoding/json"
	"github.com/RomanAgaltsev/urlcut/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/logger"
	repositoryurl "github.com/RomanAgaltsev/urlcut/internal/repository/url"
	serviceurl "github.com/RomanAgaltsev/urlcut/internal/service/url"
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

	cfg := &config.Config{
		ServerPort: "localhost:8080",
		BaseURL:    "http://localhost:8080",
		IDlength:   8,
	}

	_ = logger.Initialize()

	repo := repositoryurl.New()
	service := serviceurl.New(repo, cfg.BaseURL, cfg.IDlength)
	handlers := New(service)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.reqMethod, "/", strings.NewReader(test.reqURL))
			req.Header.Set("Content-Type", "text/plain")

			w := httptest.NewRecorder()
			h := WithLogging(handlers.Shorten)
			h(w, req)

			res := w.Result()

			assert.Equal(t, test.resStatus, res.StatusCode)
			if test.resStatus == http.StatusBadRequest {
				return
			}

			assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))

			defer func() { _ = res.Body.Close() }()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			shortenedURL := string(resBody)
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
			http.MethodGet, "", http.StatusBadRequest},
		{"[PUT] ['']",
			http.MethodPut, "", http.StatusBadRequest},
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := model.Request{
				URL: test.reqURL,
			}
			reqBytes, _ := json.Marshal(request)

			req := httptest.NewRequest(test.reqMethod, "/api/shorten", bytes.NewReader(reqBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			h := WithLogging(handlers.ShortenAPI)
			h(w, req)

			res := w.Result()
			defer func() { _ = res.Body.Close() }()

			assert.Equal(t, test.resStatus, res.StatusCode)
			if test.resStatus == http.StatusBadRequest {
				return
			}

			assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

			dec := json.NewDecoder(res.Body)
			var response model.Response
			err := dec.Decode(&response)
			require.NoError(t, err)

			shortenedURL := response.Result
			assert.Equal(t, strings.HasPrefix(shortenedURL, cfg.BaseURL), true)
		})
	}

	t.Run("[POST] [nil body]", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		h := WithLogging(handlers.ShortenAPI)
		h(w, req)

		res := w.Result()
		defer func() { _ = res.Body.Close() }()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("[POST] [Bad body]", func(t *testing.T) {
		request := struct {
			Dummy string
		}{
			Dummy: "Hi there!",
		}
		reqBytes, _ := json.Marshal(request)

		req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		h := WithLogging(handlers.ShortenAPI)
		h(w, req)

		res := w.Result()
		defer func() { _ = res.Body.Close() }()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
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
			http.MethodGet, "https://practicum.yandex.ru/", http.StatusTemporaryRedirect},
		{"[GET] [https://translate.yandex.ru/]",
			http.MethodGet, "https://translate.yandex.ru/", http.StatusTemporaryRedirect},
		{"[POST] [https://translate.yandex.ru/]",
			http.MethodPost, "https://translate.yandex.ru/", http.StatusBadRequest},
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reqPost := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.reqURL))
			reqPost.Header.Set("Content-Type", "text/plain")

			wPost := httptest.NewRecorder()
			hPost := WithLogging(handlers.Shorten)
			hPost(wPost, reqPost)

			resPost := wPost.Result()

			defer func() { _ = resPost.Body.Close() }()
			resPostBody, err := io.ReadAll(resPost.Body)
			require.NoError(t, err)

			shortenedURL := string(resPostBody)
			urlID := strings.TrimPrefix(shortenedURL, cfg.BaseURL+"/")

			req := httptest.NewRequest(test.reqMethod, "/"+urlID, nil)
			req.Header.Set("Content-Type", "text/plain")

			w := httptest.NewRecorder()
			h := WithLogging(handlers.Expand)
			h(w, req)

			res := w.Result()
			defer func() { _ = res.Body.Close() }()

			assert.Equal(t, test.resStatus, res.StatusCode)
			if res.StatusCode == http.StatusBadRequest {
				return
			}

			if assert.Contains(t, res.Header, "Location") {
				assert.Equal(t, test.reqURL, res.Header.Get("Location"))
			}
		})
	}
}
