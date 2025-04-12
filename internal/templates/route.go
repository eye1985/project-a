package templates

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"project-a/internal/middleware"
	"project-a/internal/shared"
)

func RegisterRoutes(m *middleware.Middleware, pool *pgxpool.Pool, wsUrl string, session shared.Session) {
	m.HandleFuncWithMiddleWare("GET /chat", RenderChat(&RenderChatArgs{
		WsUrl: wsUrl,
		Pool:  pool,
	}), middleware.Guard(session))
	m.HandleFunc("GET /", RenderRegisterUser)
}
