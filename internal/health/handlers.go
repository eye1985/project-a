package health

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type Handler struct {
	Pool *pgxpool.Pool
}

func NewHealthHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{pool}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	var status string
	err := h.Pool.Ping(context.Background())
	if err != nil {
		status = "No database connection established"
	}

	status = "Database: OK"

	_, _ = w.Write([]byte(status))
}
