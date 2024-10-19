package url

import (
	"net/http"
	"strings"
)

func (h *Handlers) ExpandURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urlID := strings.TrimPrefix(r.URL.Path, "/")

	expandedURL, err := h.service.ExpandURL(urlID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Location", expandedURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
