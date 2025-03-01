package health

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type IHealthController interface {
	GetHealth(mux *http.ServeMux)
}

type HealthController struct {
	pool *pgxpool.Pool
}

func NewHealthController(pool *pgxpool.Pool) *HealthController {
	return &HealthController{pool}
}

func (h *HealthController) GetHealth(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		var status string
		err := h.pool.Ping(context.Background())
		if err != nil {
			status = "No database connection established"
		}

		status = "Database: OK"

		w.Write([]byte(status))
	})
}
