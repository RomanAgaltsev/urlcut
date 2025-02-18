package url

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/RomanAgaltsev/urlcut/internal/api/middleware"
	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/logger"
	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/pkg/random"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
	"github.com/RomanAgaltsev/urlcut/internal/services"
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
		SecretKey:       "secret",
		IDlength:        idLength,
	}

	repo := repository.NewInMemoryRepository("storage.json")
	service, err := services.NewShortener(repo, cfg)
	require.NoError(t, err)
	router := chi.NewRouter()
	handlers := NewHandlers(service, cfg)

	tokenAuth := jwtauth.New("HS256", []byte(cfg.SecretKey), nil)
	router.Use(middleware.WithAuth(tokenAuth))

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

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc := resty.New().SetCookieJar(jar)

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
			req := httpc.R()
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

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc := resty.New().SetCookieJar(jar)

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
			req := httpc.R()
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
		req := httpc.R()
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

		req := httpc.R()
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

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc := resty.New().SetCookieJar(jar)

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

		req := httpc.R()
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

		for urlF, found := range reqFinds {
			assert.Truef(t, found, "url [%s] not found in response", urlF)
		}

	})

	t.Run("[POST] [nil body]", func(t *testing.T) {
		req := httpc.R()
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

		req := httpc.R()
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

		req := httpc.R()
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

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc := resty.New().SetCookieJar(jar)

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
			reqPost := httpc.R()
			reqPost.Method = http.MethodPost
			reqPost.URL = httpSrv.URL

			resPost, err := reqPost.
				SetHeader("Content-Type", ContentTypeText).
				SetBody(test.reqURL).
				Send()
			assert.NoError(t, err)

			shortenedURL := string(resPost.Body())
			urlID := strings.TrimPrefix(shortenedURL, hlp.cfg.BaseURL+"/")

			req := httpc.R()
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

func TestPingHandler(t *testing.T) {
	hlp := newHelper(t)
	hlp.router.Get("/ping", hlp.handlers.Ping)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc := resty.New().SetCookieJar(jar)

	res, err := httpc.R().
		Get(httpSrv.URL + "/ping")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode())

}

func TestUserUrlsHandler(t *testing.T) {
	hlp := newHelper(t)
	hlp.router.Get("/api/user/urls", hlp.handlers.UserUrls)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc := resty.New().SetCookieJar(jar)

	req := httpc.R()
	req.Method = http.MethodGet
	req.URL = httpSrv.URL + "/api/user/urls"

	res, err := req.
		SetHeader("Content-Type", ContentTypeText).
		Send()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, res.StatusCode())
}

func TestUserUrlsDeleteHandler(t *testing.T) {
	hlp := newHelper(t)
	hlp.router.Delete("/api/user/urls", hlp.handlers.UserUrlsDelete)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc := resty.New().SetCookieJar(jar)

	ids := []string{random.String(8), random.String(8), random.String(8)}
	idsBytes, _ := json.Marshal(ids)

	req := httpc.R()
	req.Method = http.MethodDelete
	req.URL = httpSrv.URL + "/api/user/urls"

	res, err := req.
		SetHeader("Content-Type", ContentTypeJSON).
		SetBody(bytes.NewReader(idsBytes)).
		Send()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusAccepted, res.StatusCode())
}

func TestLoggerMiddleWare(t *testing.T) {
	err := logger.Initialize()
	require.NoError(t, err)

	hlp := newHelper(t)
	hlp.router.Use(middleware.WithLogging)
	hlp.router.Post("/", hlp.handlers.Shorten)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc := resty.New().SetCookieJar(jar)

	t.Run("[POST] [LoggerMiddleware] [https://practicum.yandex.ru/]", func(t *testing.T) {
		res, err := httpc.R().
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

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc := resty.New().SetCookieJar(jar)

	t.Run("[POST] [CompressMiddleware gzip/''] [https://practicum.yandex.ru/]", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		gzipWriter := gzip.NewWriter(buf)
		_, err := gzipWriter.Write([]byte("https://practicum.yandex.ru/"))
		require.NoError(t, err)
		err = gzipWriter.Close()
		require.NoError(t, err)

		res, err := httpc.R().
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

func TestAuthIDMiddleware(t *testing.T) {
	hlp := newHelper(t)
	tokenAuth := jwtauth.New("HS256", []byte(hlp.cfg.SecretKey), nil)
	hlp.router.Use(middleware.WithID(tokenAuth))
	hlp.router.Get("/auth", hlp.handlers.UserUrls)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	httpc := resty.New()

	req := httpc.R()
	req.Method = http.MethodGet
	req.URL = httpSrv.URL + "/auth"

	res, err := req.
		SetHeader("Content-Type", ContentTypeText).
		Send()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode())

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc.SetCookieJar(jar)

	req = httpc.R()
	req.Method = http.MethodGet
	req.URL = httpSrv.URL + "/auth"

	res, err = req.
		SetHeader("Content-Type", ContentTypeText).
		Send()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, res.StatusCode())
}

func BenchmarkHandlers_Shorten(b *testing.B) {
	hlp := newHelper(&testing.T{})
	hlp.router.Post("/", hlp.handlers.Shorten)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc := resty.New().SetCookieJar(jar)

	req := httpc.R()
	req.Method = http.MethodPost
	req.URL = httpSrv.URL
	req.SetHeader("Content-Type", ContentTypeText).SetBody("https://practicum.yandex.ru/")

	b.ResetTimer()
	for range b.N {
		_, _ = req.Send()
	}
}

func BenchmarkHandlers_ShortenAPI(b *testing.B) {
	hlp := newHelper(&testing.T{})
	hlp.router.Post("/api/shorten", hlp.handlers.ShortenAPI)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc := resty.New().SetCookieJar(jar)

	req := httpc.R()
	req.Method = http.MethodPost
	req.URL = httpSrv.URL + "/api/shorten"

	req.SetHeader("Content-Type", ContentTypeJSON).SetBody(`{"url":"https://practicum.yandex.ru/"}`)

	b.ResetTimer()
	for range b.N {
		_, _ = req.Send()
	}
}

func BenchmarkHandlers_ShortenAPIBatch(b *testing.B) {
	hlp := newHelper(&testing.T{})
	hlp.router.Post("/api/shorten/batch", hlp.handlers.ShortenAPIBatch)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc := resty.New().SetCookieJar(jar)

	req := httpc.R()
	req.Method = http.MethodPost
	req.URL = httpSrv.URL + "/api/shorten/batch"

	req.SetHeader("Content-Type", ContentTypeJSON).
		SetBody(`[{"correlation_id":"kj24njF2","original_url":"https://practicum.yandex.ru/"},{"correlation_id":"87sdFin3","original_url":"https://translate.yandex.ru/"}]`)

	b.ResetTimer()
	for range b.N {
		_, _ = req.Send()
	}
}

func BenchmarkHandlers_Expand(b *testing.B) {
	hlp := newHelper(&testing.T{})
	hlp.router.Post("/", hlp.handlers.Shorten)
	hlp.router.Get("/{id}", hlp.handlers.Expand)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc := resty.New().SetCookieJar(jar)

	reqPost := httpc.R()
	reqPost.Method = http.MethodPost
	reqPost.URL = httpSrv.URL

	reqPost.SetHeader("Content-Type", ContentTypeText).SetBody("https://practicum.yandex.ru/")

	b.ResetTimer()
	for range b.N {
		b.StopTimer()
		resPost, err := reqPost.Send()
		if err != nil {
			continue
		}

		shortenedURL := string(resPost.Body())

		urlID := strings.TrimPrefix(shortenedURL, hlp.cfg.BaseURL+"/")

		req := httpc.R()
		req.Method = http.MethodGet
		req.URL = httpSrv.URL + "/" + urlID

		req.SetHeader("Content-Type", ContentTypeText)

		b.StartTimer()
		_, _ = req.Send()
	}
}

func BenchmarkHandlers_UserUrls(b *testing.B) {
	hlp := newHelper(&testing.T{})
	hlp.router.Get("/api/user/urls", hlp.handlers.UserUrls)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc := resty.New().SetCookieJar(jar)

	req := httpc.R()
	req.Method = http.MethodGet
	req.URL = httpSrv.URL + "/api/user/urls"

	req.SetHeader("Content-Type", ContentTypeText)

	b.ResetTimer()
	for range b.N {
		_, _ = req.Send()
	}
}

func BenchmarkHandlers_UserUrlsDelete(b *testing.B) {
	hlp := newHelper(&testing.T{})
	hlp.router.Delete("/api/user/urls", hlp.handlers.UserUrlsDelete)

	httpSrv := httptest.NewServer(hlp.router)
	defer httpSrv.Close()

	u, _ := url.Parse(httpSrv.URL)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{
		Name:  "jwt",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3NjU3YTI5OC0xNjMyLTQzMTUtYjc3Yi01N2QwYTFmYTFlYjQifQ.__0hZzB7EPGqGGR3o9xYsOx5ucWazs3ExB4pQ5bzjmw",
		Path:  "/",
	}})

	httpc := resty.New().SetCookieJar(jar)

	req := httpc.R()
	req.Method = http.MethodDelete
	req.URL = httpSrv.URL + "/api/user/urls"

	req.SetHeader("Content-Type", ContentTypeJSON).SetBody(`["6qxTVvsy", "RTfd56hn", "Jlfd67ds"] `)

	b.ResetTimer()
	for range b.N {
		_, _ = req.Send()
	}
}
