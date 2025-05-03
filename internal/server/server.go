package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"project-a/internal/auth"
	"project-a/internal/contacts"
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

	midWare := middleware.NewMux()
	midWare.Add(middleware.Logger)
	midWare.Add(middleware.BodyCloser)

	hub := socket.NewHub()
	go hub.Run()

	// repos
	authRepo := auth.NewRepo(pool)
	userRepo := user.NewUserRepo(pool)
	contactsRepo := contacts.NewRepo(pool)

	// services
	authService := auth.NewService(authRepo, hashKey, blockKey)

	// handlers
	healthHandler := health.NewHandler(pool)
	userHandler := user.NewUserHandler(userRepo, hub)
	contactsHandler := contacts.NewHandler(contactsRepo, userRepo)
	authHandler := auth.NewHandler(authService, authRepo, userRepo, contactsRepo)
	templateHandler := templates.NewHandler(userRepo, contactsRepo, authService, wsUrl)

	// routes
	midWare.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))
	midWare.Handle("GET /styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("web/styles"))))
	midWare.HandleFunc(
		"GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	)

	health.RegisterRoutes(midWare, healthHandler)
	auth.RegisterRoutes(midWare, authHandler, authService)
	user.RegisterRoutes(midWare, userHandler, authService)
	contacts.RegisterRoutes(midWare, contactsHandler, authService)
	socket.RegisterRoutes(midWare, hub, authService, userRepo, contactsRepo)
	templates.RegisterRoutes(midWare, templateHandler, authService)

	return http.ListenAndServe(PORT, midWare.Mux)
}
