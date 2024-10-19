package url

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *Handlers) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, _ := io.ReadAll(r.Body)
	defer func() { _ = r.Body.Close() }()

	if len(url) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//return fmt.Sprintf("%s/%s", s.baseURL, id), nil
	shortenedURL, err := h.service.ShortenURL(string(url))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(shortenedURL)))
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(shortenedURL))
	if err != nil {
		slog.Info(
			"failed to write shorten URL to response",
			"error", err.Error())
	}
}
