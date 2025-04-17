package templates

import (
	"fmt"
	"html/template"
	"net/http"
	"project-a/internal/shared"
)

const (
	path     = "web/templates"
	register = "register.gohtml"
	chat     = "chat.gohtml"
)

func RenderRegisterUser(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/%s", path, register))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RenderChat(props *RenderChatArgs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(fmt.Sprintf("%s/%s", path, chat))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sessionID := r.Context().Value(shared.SessionCtxKey).([]byte)
		user, err := props.ur.GetUserFromSessionId(string(sessionID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, &PageData{
			WsUrl:    props.wsUrl,
			Username: user.Username,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
