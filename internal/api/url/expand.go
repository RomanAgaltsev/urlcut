package url

import (
	"net/http"
	"strings"
)

func (h *Handlers) Expand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urlID := strings.TrimPrefix(r.URL.Path, "/")

	url, err := h.service.Expand(urlID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Location", url.Long)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
