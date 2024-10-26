package url

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/urlcut/internal/api/middleware"
	"github.com/RomanAgaltsev/urlcut/internal/logger"
	"github.com/RomanAgaltsev/urlcut/internal/model"
	repositoryurl "github.com/RomanAgaltsev/urlcut/internal/repository/url"
	serviceurl "github.com/RomanAgaltsev/urlcut/internal/service/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type helper struct {
	baseURL  string
	idLength int
	repo     *repositoryurl.InMemoryRepository
	service  *serviceurl.ShortenerService
	router   *chi.Mux
	handlers *Handlers
}

func newHelper() *helper {
	const (
		baseURL  = "http://localhost:8080"
		idLength = 8
	)

	repo := repositoryurl.New("storage.json")
	service := serviceurl.New(repo, baseURL, idLength)
	router := chi.NewRouter()
	handlers := New(service)

	return &helper{
		baseURL:  baseURL,
		idLength: idLength,
		repo:     repo,
		service:  service,
		router:   router,
		handlers: handlers,
	}
}

func TestShortenHandler(t *testing.T) {
	hlp := newHelper()
	hlp.router.Post("/", hlp.handlers.Shorten)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	tests := []struct {
		name      string
		reqMethod string
		reqURL    string
		resStatus int
	}{
		{"[POST] [https://practicum.yandex.ru/]", http.MethodPost, "https://practicum.yandex.ru/", http.StatusCreated},
		{"[POST] [https://translate.yandex.ru/]", http.MethodPost, "https://practicum.yandex.ru/", http.StatusCreated},
		{"[POST] ['']", http.MethodPost, "", http.StatusBadRequest},
		{"[GET] ['']", http.MethodGet, "", http.StatusMethodNotAllowed},
		{"[PUT] ['']", http.MethodPut, "", http.StatusMethodNotAllowed},
	}

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

			shortenedURL := string(res.Body())

			assert.Equal(t, "text/plain", res.Header().Get("Content-Type"))
			assert.Equal(t, strings.HasPrefix(shortenedURL, hlp.baseURL), true)
		})
	}
}

func TestShortenAPIHandler(t *testing.T) {
	hlp := newHelper()
	hlp.router.Post("/api/shorten", hlp.handlers.ShortenAPI)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	tests := []struct {
		name      string
		reqMethod string
		reqURL    string
		resStatus int
	}{
		{"[POST] [https://practicum.yandex.ru/]", http.MethodPost, "https://practicum.yandex.ru/", http.StatusCreated},
		{"[POST] [https://translate.yandex.ru/]", http.MethodPost, "https://practicum.yandex.ru/", http.StatusCreated},
		{"[POST] ['']", http.MethodPost, "", http.StatusBadRequest},
		{"[GET] ['']", http.MethodGet, "", http.StatusMethodNotAllowed},
		{"[PUT] ['']", http.MethodPut, "", http.StatusMethodNotAllowed},
	}

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
			assert.Equal(t, strings.HasPrefix(shortenedURL, hlp.baseURL), true)
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
	hlp := newHelper()
	hlp.router.Post("/", hlp.handlers.Shorten)
	hlp.router.Get("/{id}", hlp.handlers.Expand)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	tests := []struct {
		name      string
		reqMethod string
		reqURL    string
		resStatus int
	}{
		{"[GET] [https://practicum.yandex.ru/]", http.MethodGet, "https://practicum.yandex.ru/", http.StatusOK},
		{"[GET] [https://translate.yandex.ru/]", http.MethodGet, "https://translate.yandex.ru/", http.StatusOK},
		{"[POST] [https://translate.yandex.ru/]", http.MethodPost, "https://translate.yandex.ru/", http.StatusMethodNotAllowed},
	}

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
			urlID := strings.TrimPrefix(shortenedURL, hlp.baseURL+"/")

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

func TestLoggerMiddleWare(t *testing.T) {
	err := logger.Initialize()
	require.NoError(t, err)

	hlp := newHelper()
	hlp.router.Use(middleware.WithLogging)
	hlp.router.Post("/", hlp.handlers.Shorten)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	t.Run("[POST] [LoggerMiddleware] [https://practicum.yandex.ru/]", func(t *testing.T) {
		res, err := resty.
			New().
			R().
			SetHeader("Content-Type", "text/plain").
			SetBody("https://practicum.yandex.ru/").
			Post(httpSrv.URL)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, res.StatusCode())
		if res.StatusCode() == http.StatusBadRequest || res.StatusCode() == http.StatusMethodNotAllowed {
			return
		}

		shortenedURL := string(res.Body())

		assert.Equal(t, "text/plain", res.Header().Get("Content-Type"))
		assert.Equal(t, strings.HasPrefix(shortenedURL, hlp.baseURL), true)
	})
}

func TestCompressMiddleware(t *testing.T) {
	hlp := newHelper()
	hlp.router.Use(middleware.WithGzip)
	hlp.router.Post("/compress", hlp.handlers.Shorten)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	t.Run("[POST] [CompressMiddleware gzip/''] [https://practicum.yandex.ru/]", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		gzipWriter := gzip.NewWriter(buf)
		_, err := gzipWriter.Write([]byte("https://practicum.yandex.ru/"))
		require.NoError(t, err)
		err = gzipWriter.Close()
		require.NoError(t, err)

		res, err := resty.
			New().
			R().
			SetHeader("Content-Encoding", "gzip").
			SetHeader("Accept-Encoding", "").
			SetBody(buf.Bytes()).
			Post(httpSrv.URL + "/compress")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, res.StatusCode())
		if res.StatusCode() == http.StatusBadRequest || res.StatusCode() == http.StatusMethodNotAllowed {
			return
		}

		shortenedURL := string(res.Body())

		assert.Equal(t, "text/plain", res.Header().Get("Content-Type"))
		assert.Equal(t, strings.HasPrefix(shortenedURL, hlp.baseURL), true)
	})

	t.Run("[POST] [CompressMiddleware ''/gzip] [https://practicum.yandex.ru/]", func(t *testing.T) {
		res, err := resty.
			New().
			R().
			SetHeader("Accept-Encoding", "gzip").
			SetBody("https://practicum.yandex.ru/").
			Post(httpSrv.URL + "/compress")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, res.StatusCode())
		if res.StatusCode() == http.StatusBadRequest || res.StatusCode() == http.StatusMethodNotAllowed {
			return
		}

		shortenedURL := string(res.Body())

		assert.Equal(t, "text/plain", res.Header().Get("Content-Type"))
		assert.Equal(t, strings.HasPrefix(shortenedURL, hlp.baseURL), true)
	})
}
