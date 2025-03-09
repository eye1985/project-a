package server

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"html/template"
	"log"
	"net/http"
	"project-a/modules/health"
	"project-a/modules/users"
	"project-a/server/middleware"
	"project-a/socket"
)

const PORT = ":8080"

type JSON struct {
	Message string `json:"message"`
}

type Person struct {
	Username string
	Email    string
}

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func Serve(pool *pgxpool.Pool) error {
	midWare := middleware.NewMiddlewareMux()
	mux := midWare.Mux
	midWare.Add(LoggerMiddleware)

	// Chat
	wsHandler := socket.ServeWs()

	// Modules
	userModule := users.NewUserModule(pool, mux)

	// Controllers
	userModule.UserController.GetUsers()
	userModule.UserController.RegisterUser()
	health.NewHealthController(pool).GetHealth(midWare)

	midWare.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	midWare.HandleFunc("GET /ws", wsHandler)
	midWare.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	midWare.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.gohtml")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rows, err := pool.Query(context.Background(), "SELECT username, email FROM users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		persons := []Person{}

		for rows.Next() {
			var username, email string

			if err := rows.Scan(&username, &email); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			persons = append(persons, Person{
				username,
				email,
			})
		}

		if err := tmpl.Execute(w, persons); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	return http.ListenAndServe(PORT, mux)
}
