package url

import (
	"encoding/json"
	"github.com/RomanAgaltsev/urlcut/internal/model"
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

func (h *Handlers) ShortenAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(r.Body)

	var req model.Request
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Info(
			"failed to unmarshal long URL",
			"error", err.Error())
		return
	}

	longURL := req.URL
	if len(longURL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := h.service.Shorten(longURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Info(
			"failed to short URL",
			"error", err.Error())
		return
	}

	res, err := json.Marshal(model.Response{Result: url.Short()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Info(
			"failed to marshal shorten URL",
			"error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(res)))
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Info(
			"failed to write shorten URL to response",
			"error", err.Error())
	}
}
