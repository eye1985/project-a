package server

import (
	"context"
	"github.com/jackc/pgx/v5"
	"html/template"
	"net/http"
)

const PORT = ":8080"

type JSON struct {
	Message string `json:"message"`
}

type Person struct {
	Username string
	Email    string
}

func Serve(pgUrl *string) error {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pgConn, err := pgx.Connect(context.Background(), *pgUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer pgConn.Close(context.Background())

		rows, err := pgConn.Query(context.Background(), "SELECT username, email FROM users")
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

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		pgConn, err := pgx.Connect(context.Background(), *pgUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer pgConn.Close(context.Background())

		var status string
		err = pgConn.Ping(context.Background())
		if err != nil {
			status = "No database connection established"
		}

		status = "Database: OK"

		w.Write([]byte(status))
	})

	return http.ListenAndServe(PORT, mux)
}
