package user

import (
	"encoding/json"
	"log"
	"net/http"
	"project-a/internal/shared"
	"project-a/internal/socket"
	"strings"
)

type Handler struct {
	Repo shared.UserRepository
	Hub  *socket.Hub
}

func (h *Handler) GetUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	users, err := h.Repo.GetUsers()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid request content-type", http.StatusUnsupportedMediaType)
		return
	}

	user := &shared.User{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "Missing required field: username", http.StatusBadRequest)
		return
	}
	if user.Email == "" {
		http.Error(w, "Missing required field: email", http.StatusBadRequest)
		return
	}

	_, err = h.Repo.InsertUser(user.Username, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) UpdateUserName(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "Invalid request content-type", http.StatusUnsupportedMediaType)
		return
	}

	sessionID := r.Context().Value(shared.SessionCtxKey).([]byte)
	u, err := h.Repo.GetUserFromSessionId(string(sessionID))
	if err != nil {
		http.Error(w, "no session", http.StatusInternalServerError)
		return
	}

	var newUser shared.User
	err = json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(newUser.Username) < 3 {
		http.Error(w, "Username must be at least 3 characters long", http.StatusBadRequest)
		return
	}

	err = h.Repo.UpdateUserName(newUser.Username, u.Id)
	if err != nil {
		log.Printf("%s", err.Error())
		http.Error(w, "could not update username", http.StatusInternalServerError)
		return
	}

	h.Hub.UpdateNameChannel(u.Id, newUser.Username)

	w.WriteHeader(http.StatusNoContent)
}

func NewUserHandler(repo shared.UserRepository, hub *socket.Hub) *Handler {
	return &Handler{
		Repo: repo,
		Hub:  hub,
	}
}
