package url

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *Handlers) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	longURL, _ := io.ReadAll(r.Body)
	defer func() { _ = r.Body.Close() }()

	if len(longURL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := h.service.Shorten(string(longURL))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortURL := url.Short()

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(shortURL)))
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(shortURL))
	if err != nil {
		slog.Info(
			"failed to write shorten URL to response",
			"error", err.Error())
	}
}
