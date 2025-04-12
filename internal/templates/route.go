package templates

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"project-a/internal/middleware"
)

func RegisterRoutes(m *middleware.Middleware, pool *pgxpool.Pool, wsUrl string) {
	m.HandleFunc("GET /", RenderChat(&RenderChatArgs{
		WsUrl: wsUrl,
		Pool:  pool,
	}))
}
