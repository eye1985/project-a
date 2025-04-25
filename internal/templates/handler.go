package templates

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"project-a/internal/shared"
)

const (
	path       = "web/templates"
	baseLayout = "base-layout.gohtml"
	register   = "register.gohtml"
	chat       = "chat.gohtml"
	profile    = "profile.gohtml"
)

type Handler struct {
	userRepo    shared.UserRepository
	authService shared.Session
	wsUrl       string
}

func (h *Handler) RenderRegisterUser(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		fmt.Sprintf("%s/%s", path, register),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RenderChat(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		fmt.Sprintf("%s/%s", path, baseLayout),
		fmt.Sprintf("%s/%s", path, chat),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID := r.Context().Value(shared.SessionCtxKey).([]byte)
	u, err := h.userRepo.GetUserFromSessionId(string(sessionID))
	if err != nil {
		http.Error(w, "no session", http.StatusInternalServerError)
		return
	}

	log.Printf("wsUrl: %v", h.wsUrl)
	log.Printf("username: %v", u.Username)

	if err := tmpl.Execute(w, &PageData{
		WsUrl:    h.wsUrl,
		Username: u.Username,
		Title:    "Chat",
		Css:      "chat.css",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RenderProfile(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		fmt.Sprintf("%s/%s", path, baseLayout),
		fmt.Sprintf("%s/%s", path, profile),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sessionID := r.Context().Value(shared.SessionCtxKey).([]byte)
	u, err := h.userRepo.GetUserFromSessionId(string(sessionID))
	if err != nil {
		http.Error(w, "no session", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, &PageData{
		Username: u.Username,
		Title:    "Profile",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewHandler(userRepo shared.UserRepository, authService shared.Session, wsUrl string) *Handler {
	return &Handler{
		userRepo:    userRepo,
		authService: authService,
		wsUrl:       wsUrl,
	}
}
