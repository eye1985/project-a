package server

import (
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"os"
	"project-a/internal/auth"
	"project-a/internal/contacts"
	"project-a/internal/email"
	"project-a/internal/health"
	"project-a/internal/middleware"
	"project-a/internal/socket"
	"project-a/internal/templates"
	"project-a/internal/user"
	"project-a/internal/util"
)

const PORT = ":3000"

type ServeArgs struct {
	HashKey       string
	BlockKey      string
	WsUrl         string
	Origin        string
	MailSendToken string
}

func Serve(pool *pgxpool.Pool, args *ServeArgs) error {
	hashKey := args.HashKey
	blockKey := args.BlockKey
	wsUrl := args.WsUrl
	origin := args.Origin
	mailSendToken := args.MailSendToken

	isDev, _ := os.LookupEnv("IS_DEV")
	dev := false

	if isDev == "true" {
		dev = true
	}

	var csrfHandler func(http.Handler) http.Handler
	if dev {
		csrfHandler = csrf.Protect(
			util.GetCSRFKey(),
			csrf.SameSite(csrf.SameSiteStrictMode),
			csrf.Secure(false),
			csrf.TrustedOrigins([]string{"localhost:3000"}),
			csrf.Path("/"),
		)
	} else {
		csrfHandler = csrf.Protect(
			util.GetCSRFKey(),
			csrf.SameSite(csrf.SameSiteStrictMode),
			csrf.Secure(false),
			csrf.Path("/"),
		)
	}

	midWare := middleware.NewMux()
	midWare.Add(middleware.Logger)
	midWare.Add(middleware.BodyCloser)
	midWare.Add(middleware.NoCache)

	hub := socket.NewHub()
	go hub.Run()

	// repos
	authRepo := auth.NewRepo(pool)
	userRepo := user.NewUserRepo(pool)
	contactsRepo := contacts.NewRepo(pool)
	emailRepo := email.NewRepo(pool)

	// services
	authService := auth.NewService(authRepo, hashKey, blockKey)

	// handlers
	healthHandler := health.NewHandler(pool)
	userHandler := user.NewUserHandler(userRepo, hub)
	contactsHandler := contacts.NewHandler(contactsRepo, userRepo)
	authHandler := auth.NewHandler(
		&auth.NewHandlerArgs{
			AuthService:   authService,
			Repo:          authRepo,
			UserRepo:      userRepo,
			ContactsRepo:  contactsRepo,
			EmailRepo:     emailRepo,
			MailSendToken: mailSendToken,
			Origin:        origin,
		},
	)
	templateHandler := templates.NewHandler(userRepo, contactsRepo, authService, wsUrl, dev)

	// routes
	midWare.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))
	midWare.HandleFunc("GET /favicon.ico", http.NotFound)
	midWare.HandleFunc("GET /.well-known/", http.NotFound)

	health.RegisterRoutes(midWare, healthHandler)
	auth.RegisterRoutes(
		midWare,
		authHandler,
		authService,
		csrfHandler,
	)
	user.RegisterRoutes(
		midWare,
		userHandler,
		authService,
	)
	contacts.RegisterRoutes(
		midWare,
		contactsHandler,
		authService,
	)
	socket.RegisterRoutes(
		midWare,
		hub,
		authService,
		userRepo,
		contactsRepo,
		origin,
	)
	templates.RegisterRoutes(
		midWare,
		templateHandler,
		authService,
		csrfHandler,
	)

	return http.ListenAndServe(PORT, midWare.Mux)
}
