package url

import (
	"log/slog"
	"net/http"
	"time"

	servicesurl "github.com/RomanAgaltsev/urlcut/internal/services/url"
)

type Handlers struct {
	service servicesurl.Service
}

func NewHandlers(service servicesurl.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func WithLogging(h http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		respData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   respData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		slog.Info(
			"got request",
			slog.String("uri", r.RequestURI),
			slog.String("method", r.Method),
			slog.Int("status", respData.status),
			slog.Duration("duration", duration),
			slog.Int("size", respData.size),
		)
	}
	return logFn
}
