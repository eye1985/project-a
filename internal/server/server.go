package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"project-a/internal/auth"
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
	hashKey, ok := os.LookupEnv("HASH_KEY")
	if !ok {
		log.Fatalf("HASH_KEY environment variable not set")
	}
	blockKey, ok := os.LookupEnv("BLOCK_KEY")
	if !ok {
		log.Fatalf("BLOCK_KEY environment variable not set")
	}

	midWare := middleware.NewMiddlewareMux()
	midWare.Add(middleware.Logger)
	mux := midWare.Mux

	hub := socket.NewHub()
	go hub.Run()

	// repos
	authRepo := auth.NewAuthRepo(pool)

	// services
	userService := user.NewUserService(pool)
	authService := auth.NewAuthService(authRepo, hashKey, blockKey)

	// handlers
	healthHandler := health.NewHealthHandler(pool)
	userHandler := user.NewUserHandler(pool)
	authHandler := auth.NewAuthHandler(authService, authRepo, userService)

	// routes
	midWare.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))
	midWare.Handle("GET /styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("web/styles"))))
	midWare.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	health.RegisterRoutes(midWare, healthHandler)
	auth.RegisterRoutes(midWare, authHandler, authService)
	user.RegisterRoutes(midWare, userHandler)
	socket.RegisterRoutes(midWare, hub)
	templates.RegisterRoutes(&templates.RegisterRoutesArgs{
		Middleware:  midWare,
		WsUrl:       wsUrl,
		Session:     authService,
		UserService: userService,
	})

	return http.ListenAndServe(PORT, mux)
}
