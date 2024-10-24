package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	writer   http.ResponseWriter
	gzwriter *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		writer:   w,
		gzwriter: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.writer.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.gzwriter.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	c.writer.Header().Set("Content-Encoding", "gzip")
	c.writer.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.gzwriter.Close()
}

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

func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.gzreader.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}
	return c.gzreader.Close()
}

func WithGzip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		tempWriter := w
		//		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		//			compWriter := newCompressWriter(w)
		//			tempWriter = compWriter
		//			defer func() { _ = compWriter.Close() }()
		//		}

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
