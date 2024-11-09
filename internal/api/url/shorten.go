package url

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/RomanAgaltsev/urlcut/internal/model"
)

func (h *Handlers) Shorten(w http.ResponseWriter, r *http.Request) {
	longURL, _ := io.ReadAll(r.Body)
	defer func() { _ = r.Body.Close() }()

	if len(longURL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := h.shortener.Shorten(string(longURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortURL := url.Short()

	w.Header().Set("Content-Type", ContentTypeText)
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
	dec := json.NewDecoder(r.Body)

	var req model.Request
	if err := dec.Decode(&req); err != nil {
		slog.Info(
			"failed to unmarshal long URL",
			"error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	longURL := req.URL
	if len(longURL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := h.shortener.Shorten(longURL)
	if err != nil {
		slog.Info(
			"failed to short URL",
			"error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(model.Response{Result: url.Short()})
	if err != nil {
		slog.Info(
			"failed to marshal shorten URL",
			"error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
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

func (h *Handlers) ShortenAPIBatch(w http.ResponseWriter, r *http.Request) {
	batch := make([]model.BatchRequest, 0)

	dec := json.NewDecoder(r.Body)

	_, err := dec.Token()
	if err != nil {
		slog.Info(
			"failed to decode batch",
			"error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for dec.More() {
		var batchReq model.BatchRequest
		if err := dec.Decode(&batchReq); err != nil {
			slog.Info(
				"failed to decode batch element",
				"error", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		batch = append(batch, batchReq)
	}

	if len(batch) == 0 {
		slog.Info("got empty batch")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	batchShortened, err := h.shortener.ShortenBatch(batch)
	if err != nil {
		slog.Info(
			"failed to short URL",
			"error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	//	for _, batchelem := range batchShortened {
	//		enc.Encode(batchelem)
	//	}
	err = enc.Encode(batchShortened)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Info(
			"failed to encode batch",
			"error", err.Error())
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	//w.Header().Set("Content-Length", strconv.Itoa(len(res)))
	w.WriteHeader(http.StatusCreated)

	//	_, err = w.Write(res)
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		slog.Info(
	//			"failed to write shorten URL to response",
	//			"error", err.Error())
	//	}
}
