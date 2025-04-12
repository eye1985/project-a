package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"project-a/internal/health"
	"project-a/internal/middleware"
	"project-a/internal/socket"
	"project-a/internal/templates"
	"project-a/internal/user"
)

const PORT = ":3000"

func Serve(pool *pgxpool.Pool) error {
	wsUrl, ok := os.LookupEnv("WS_URL")
	if !ok {
		log.Fatalf("WS_URL environment variable not set")
	}

	midWare := middleware.NewMiddlewareMux()
	midWare.Add(middleware.Logger)
	mux := midWare.Mux

	hub := socket.NewHub()
	go hub.Run()

	midWare.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))
	midWare.Handle("GET /styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("web/styles"))))
	midWare.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	healthHandler := health.NewHealthHandler(pool)
	userHandler := user.NewUserHandler(pool)

	health.RegisterRoutes(midWare, healthHandler)
	user.RegisterRoutes(midWare, userHandler)
	socket.RegisterRoutes(midWare, hub)
	templates.RegisterRoutes(midWare, pool, wsUrl)

	return http.ListenAndServe(PORT, mux)
}
