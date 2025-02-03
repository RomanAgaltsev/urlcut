package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressReader struct {
	reader   io.ReadCloser
	gzreader *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		reader:   r,
		gzreader: zr,
	}, nil
}

// Read читает сжатые данные.
func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.gzreader.Read(p)
}

// Close закрывает чтение.
func (c *compressReader) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}
	return c.gzreader.Close()
}

// WithGzip выполняет роль миддлваре сжатия данных.
func WithGzip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			compReader, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			defer func() { _ = compReader.Close() }()
			r.Body = compReader
		}
		h.ServeHTTP(w, r)
	})
}
