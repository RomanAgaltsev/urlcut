package url

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RomanAgaltsev/urlcut/internal/api/middleware"
	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/logger"
	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
	"github.com/RomanAgaltsev/urlcut/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type helper struct {
	cfg      *config.Config
	router   *chi.Mux
	handlers *Handlers

	shortener  interfaces.Service
	repository interfaces.Repository
}

func newHelper(t *testing.T) *helper {
	const (
		serverPort = "localhost:8080"
		baseURL    = "http://localhost:8080"
		idLength   = 8
	)

	cfg := &config.Config{
		ServerPort:      serverPort,
		BaseURL:         baseURL,
		FileStoragePath: "",
		DatabaseDSN:     "",
		IDlength:        idLength,
	}

	repo := repository.NewInMemoryRepository("storage.json")
	service, err := services.NewShortener(repo, cfg)
	require.NoError(t, err)
	router := chi.NewRouter()
	handlers := NewHandlers(service, cfg)

	return &helper{
		cfg:        cfg,
		repository: repo,
		shortener:  service,
		router:     router,
		handlers:   handlers,
	}
}

func TestShortenHandler(t *testing.T) {
	hlp := newHelper(t)
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
			//			req.SetCookie(&http.Cookie{
			//					Name:   "jwt",
			//					Value:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJiN2FlNTAyYS1kNDEyLTRjODQtOTE0MS05Y2M0YjEwMjU3MjgifQ.a2WtJfOI25Qizm0pD_YcmlKL7Lwr3O2BJc4XXzfNyHA",
			//					Path:   "/",
			//					//MaxAge: 300,
			//					Secure: true,
			//			})
			req.SetContext(context.WithValue(context.Background(), "uid", map[string]interface{}{"uid": "b7ae502a-d412-4c84-9141-9cc4b1025728"}))

			req.Method = test.reqMethod
			req.URL = httpSrv.URL

			res, err := req.
				SetHeader("Content-Type", ContentTypeText).
				SetBody(test.reqURL).
				Send()
			assert.NoError(t, err)

			assert.Equal(t, test.resStatus, res.StatusCode())
			if test.resStatus == http.StatusBadRequest || test.resStatus == http.StatusMethodNotAllowed {
				return
			}

			shortenedURL := string(res.Body())

			assert.Equal(t, ContentTypeText, res.Header().Get("Content-Type"))
			assert.Equal(t, strings.HasPrefix(shortenedURL, hlp.cfg.BaseURL), true)
		})
	}
}

func TestShortenAPIHandler(t *testing.T) {
	hlp := newHelper(t)
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

			request := model.URLDTO{
				URL: test.reqURL,
			}
			reqBytes, _ := json.Marshal(request)

			res, err := req.
				SetHeader("Content-Type", ContentTypeJSON).
				SetBody(bytes.NewReader(reqBytes)).
				Send()
			assert.NoError(t, err)

			assert.Equal(t, test.resStatus, res.StatusCode())
			if test.resStatus == http.StatusBadRequest || test.resStatus == http.StatusMethodNotAllowed {
				return
			}
			assert.Equal(t, ContentTypeJSON, res.Header().Get("Content-Type"))

			dec := json.NewDecoder(bytes.NewReader(res.Body()))
			var response model.ResultDTO
			err = dec.Decode(&response)
			require.NoError(t, err)

			shortenedURL := response.Result
			assert.Equal(t, strings.HasPrefix(shortenedURL, hlp.cfg.BaseURL), true)
		})
	}

	t.Run("[POST] [nil body]", func(t *testing.T) {
		req := resty.New().R()
		req.Method = http.MethodPost
		req.URL = httpSrv.URL + "/api/shorten"

		res, err := req.
			SetHeader("Content-Type", ContentTypeJSON).
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
			SetHeader("Content-Type", ContentTypeJSON).
			SetBody(bytes.NewReader(reqBytes)).
			Send()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})
}

func TestShortenAPIBatchHandler(t *testing.T) {
	hlp := newHelper(t)
	hlp.router.Post("/api/shorten/batch", hlp.handlers.ShortenAPIBatch)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	t.Run("[POST] [urls]", func(t *testing.T) {
		type request struct {
			CorrelationID string `json:"correlation_id"`
			OriginalURL   string `json:"original_url"`
		}

		type response struct {
			CorrelationID string `json:"correlation_id"`
			ShortURL      string `json:"short_url"`
		}

		requests := []request{
			{
				CorrelationID: "kj24njF2",
				OriginalURL:   "https://practicum.yandex.ru/",
			},
			{
				CorrelationID: "87sdFin3",
				OriginalURL:   "https://translate.yandex.ru/",
			},
		}

		reqFinds := map[string]bool{"https://practicum.yandex.ru/": false, "https://translate.yandex.ru/": false}

		req := resty.New().R()
		req.Method = http.MethodPost
		req.URL = httpSrv.URL + "/api/shorten/batch"

		reqBytes, _ := json.Marshal(requests)

		res, err := req.
			SetHeader("Content-Type", ContentTypeJSON).
			SetBody(bytes.NewReader(reqBytes)).
			Send()
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, res.StatusCode())
		assert.Equal(t, ContentTypeJSON, res.Header().Get("Content-Type"))

		dec := json.NewDecoder(bytes.NewReader(res.Body()))
		var responses []response
		err = dec.Decode(&responses)
		require.NoError(t, err)

		assert.Equal(t, len(requests), len(responses))

		for _, req := range requests {
			for _, res := range responses {
				if req.CorrelationID == res.CorrelationID {
					reqFinds[req.OriginalURL] = true
					assert.True(t, strings.HasPrefix(res.ShortURL, hlp.cfg.BaseURL))
				}
			}
		}

		for url, found := range reqFinds {
			assert.Truef(t, found, "url [%s] not found in response", url)
		}

	})

	t.Run("[POST] [nil body]", func(t *testing.T) {
		req := resty.New().R()
		req.Method = http.MethodPost
		req.URL = httpSrv.URL + "/api/shorten/batch"

		res, err := req.
			SetHeader("Content-Type", ContentTypeJSON).
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
		req.URL = httpSrv.URL + "/api/shorten/batch"

		res, err := req.
			SetHeader("Content-Type", ContentTypeJSON).
			SetBody(bytes.NewReader(reqBytes)).
			Send()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode())
	})

	t.Run("[POST] [Empty body]", func(t *testing.T) {
		type request struct {
			CorrelationID string `json:"correlation_id"`
			OriginalURL   string `json:"original_url"`
		}

		requests := []request{}
		reqBytes, _ := json.Marshal(requests)

		req := resty.New().R()
		req.Method = http.MethodPost
		req.URL = httpSrv.URL + "/api/shorten/batch"

		res, err := req.
			SetHeader("Content-Type", ContentTypeJSON).
			SetBody(bytes.NewReader(reqBytes)).
			Send()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})
}

func TestExpandHandler(t *testing.T) {
	hlp := newHelper(t)
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
				SetHeader("Content-Type", ContentTypeText).
				SetBody(test.reqURL).
				Send()
			assert.NoError(t, err)

			shortenedURL := string(resPost.Body())
			urlID := strings.TrimPrefix(shortenedURL, hlp.cfg.BaseURL+"/")

			req := resty.New().R()
			req.Method = test.reqMethod
			req.URL = httpSrv.URL + "/" + urlID

			res, err := req.
				SetHeader("Content-Type", ContentTypeText).
				Send()
			assert.NoError(t, err)

			assert.Equal(t, test.resStatus, res.StatusCode())
			if res.StatusCode() == http.StatusBadRequest || res.StatusCode() == http.StatusMethodNotAllowed {
				return
			}
		})
	}

	t.Run("[GET] [not found]", func(t *testing.T) {
		res, err := resty.
			New().
			R().
			Get(httpSrv.URL + "/h7Ds18sD")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode())
	})
}

//func TestPingHandler(t *testing.T) {
//	hlp := newHelper(t)
//	hlp.router.Get("/ping", hlp.handlers.Ping)
//
//	httpSrv := httptest.NewServer(hlp.router)
//	defer httpSrv.Close()
//
//	t.Run("[GET] [ping]", func(t *testing.T) {
//		res, err := resty.
//			New().
//			R().
//			Get(httpSrv.URL + "/ping")
//		assert.NoError(t, err)
//		assert.Equal(t, http.StatusOK, res.StatusCode())
//	})
//}

func TestLoggerMiddleWare(t *testing.T) {
	err := logger.Initialize()
	require.NoError(t, err)

	hlp := newHelper(t)
	hlp.router.Use(middleware.WithLogging)
	hlp.router.Post("/", hlp.handlers.Shorten)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	t.Run("[POST] [LoggerMiddleware] [https://practicum.yandex.ru/]", func(t *testing.T) {
		res, err := resty.
			New().
			R().
			SetHeader("Content-Type", ContentTypeText).
			SetBody("https://practicum.yandex.ru/").
			Post(httpSrv.URL)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, res.StatusCode())
		if res.StatusCode() == http.StatusBadRequest || res.StatusCode() == http.StatusMethodNotAllowed {
			return
		}

		shortenedURL := string(res.Body())

		assert.Equal(t, ContentTypeText, res.Header().Get("Content-Type"))
		assert.Equal(t, strings.HasPrefix(shortenedURL, hlp.cfg.BaseURL), true)
	})
}

func TestCompressMiddleware(t *testing.T) {
	hlp := newHelper(t)
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

		assert.Equal(t, ContentTypeText, res.Header().Get("Content-Type"))
		assert.Equal(t, strings.HasPrefix(shortenedURL, hlp.cfg.BaseURL), true)
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

		assert.Equal(t, ContentTypeText, res.Header().Get("Content-Type"))
		assert.Equal(t, strings.HasPrefix(shortenedURL, hlp.cfg.BaseURL), true)
	})
}
