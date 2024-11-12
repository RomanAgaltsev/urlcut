package url

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Expand выполняет обработку запроса на получение оригинального URL.
func (h *Handlers) Expand(w http.ResponseWriter, r *http.Request) {
	// Получаем идентификатор из параметров URL
	urlID := chi.URLParam(r, "id")

	// Получаем URL по идентификатору
	url, err := h.shortener.Expand(urlID)
	if err != nil {
		slog.Info(
			"failed to expand URL",
			"error", err.Error())
		http.Error(w, "please look at logs", http.StatusNotFound)
		return
	}

	// По идентификатору ничего не нашли
	if len(url.Long) == 0 {
		http.Error(w, "URL ID was not found in repository", http.StatusNotFound)
		return
	}

	// Пишем заголовки
	w.Header().Set("Location", url.Long)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
