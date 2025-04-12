package templates

import (
	"context"
	"html/template"
	"net/http"
)

func RenderChat(props *RenderChatArgs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("web/index.gohtml")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rows, err := props.Pool.Query(context.Background(), "SELECT username, email FROM users")
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

		if err := tmpl.Execute(w, &PageData{
			Person: persons,
			WsUrl:  props.WsUrl,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
