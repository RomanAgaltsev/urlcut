package url

import (
	"context"
	"net/http"

	"github.com/RomanAgaltsev/urlcut/internal/database"
)

// Ping выполняет обработку запроса на пинг хранилища.
func (h *Handlers) Ping(w http.ResponseWriter, _ *http.Request) {
	db, err := database.NewConnection(context.Background(), "pgx", h.cfg.DatabaseDSN)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() { _ = db.Close() }()

	w.WriteHeader(http.StatusOK)
}
