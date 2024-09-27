package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	type request struct {
		method      string
		contentType string
		body        string
	}
	type response struct {
		statusCode    int
		contentType   string
		contentLength int
		body          string
	}

	brResponse := struct {
		statusCode    int
		contentType   string
		contentLength int
		body          string
	}{http.StatusBadRequest, "text/plain", 11, "Bad request"}

	tests := []struct {
		name string
		req  request
		res  response
	}{
		{"[POST] [text/plain] [https://practicum.yandex.ru/]",
			request{http.MethodPost, "text/plain", "https://practicum.yandex.ru/"},
			response{http.StatusCreated, "text/plain", 30, "http://localhost:8080/EwHXdJfB"},
		},
		{"[POST] [text/plain] ['']",
			request{http.MethodPost, "text/plain", ""}, brResponse,
		},
		{"[GET] [text/plain] ['']",
			request{http.MethodGet, "text/plain", ""}, brResponse,
		},
		{"[PUT] [text/plain] ['']",
			request{http.MethodGet, "text/plain", ""}, brResponse,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.req.method, "/", strings.NewReader(test.req.body))
			req.Header.Set("Content-Type", test.req.contentType)

			w := httptest.NewRecorder()
			shortenHandler(w, req)

			res := w.Result()

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			contentLength, err := strconv.Atoi(res.Header.Get("Content-Length"))

			assert.Equal(t, test.res.statusCode, res.StatusCode)
			assert.Equal(t, test.res.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.res.contentLength, contentLength)
			assert.Equal(t, test.res.body, string(resBody))
		})
	}
}

func TestExpandHandler(t *testing.T) {
	type request struct {
		method      string
		urlID       string
		contentType string
	}
	type response struct {
		statusCode int
		header     string
		url        string
	}

	//brResponse := response{http.StatusBadRequest, "text/plain", 11, "Bad request"}

	tests := []struct {
		name string
		req  request
		res  response
	}{
		{"[GET] [EwHXdJfB] [text/plain]",
			request{http.MethodGet, "EwHXdJfB", "text/plain"},
			response{http.StatusTemporaryRedirect, "Location", "https://practicum.yandex.ru/"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.req.method, "/"+test.req.urlID, nil)
			req.Header.Set("Content-Type", test.req.contentType)

			w := httptest.NewRecorder()
			expandHandler(w, req)

			res := w.Result()

			assert.Equal(t, test.res.statusCode, res.StatusCode)
			if assert.Contains(t, res.Header, test.res.header) {
				assert.Equal(t, test.res.url, res.Header.Get(test.res.header))
			}
		})
	}
}
