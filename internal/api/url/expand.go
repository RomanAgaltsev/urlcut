package url

import (
	"log/slog"
	"net/http"
	"strings"
)

// Expand выполняет обработку запроса на получение оригинального URL.
func (h *Handlers) Expand(w http.ResponseWriter, r *http.Request) {
	// Отсекаем слэш
	urlID := strings.TrimPrefix(r.URL.Path, "/")

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
