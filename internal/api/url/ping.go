package url

import "net/http"

// Ping выполняет обработку запроса на пинг хранилища.
func (h *Handlers) Ping(w http.ResponseWriter, _ *http.Request) {
	if err := h.shortener.Check(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
