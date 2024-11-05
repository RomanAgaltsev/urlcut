package url

import "net/http"

func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	if err := h.shortener.Check(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
