package url

import (
	"log/slog"
	"net/http"
	"strings"
)

func (h *Handlers) Expand(w http.ResponseWriter, r *http.Request) {
	urlID := strings.TrimPrefix(r.URL.Path, "/")

	url, err := h.shortener.Expand(urlID)
	if err != nil {
		slog.Info(
			"failed to expand URL",
			"error", err.Error())
		http.Error(w, "please look at logs", http.StatusNotFound)
		return
	}

	if len(url.Long) == 0 {
		http.Error(w, "URL ID was not found in repository", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url.Long)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
