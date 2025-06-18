package user

import (
	"encoding/json"
	"log"
	"net/http"
	"project-a/internal/consts"
	"project-a/internal/interfaces"
	"project-a/internal/model"
	"project-a/internal/socket"
)

type Handler struct {
	Repo interfaces.UserRepository
	Hub  *socket.Hub
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	users, err := h.Repo.GetUsers(r.Context())

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
	user := &model.User{}
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

	_, err = h.Repo.InsertUser(r.Context(), user.Username, user.Email)
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
	sessionID := r.Context().Value(consts.SessionCtxKey).([]byte)
	u, err := h.Repo.GetUserFromSessionId(r.Context(), string(sessionID))
	if err != nil {
		http.Error(w, "no session", http.StatusInternalServerError)
		return
	}

	var newUser model.User
	err = json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(newUser.Username) < 3 {
		http.Error(w, "Username must be at least 3 characters long", http.StatusBadRequest)
		return
	}

	err = h.Repo.UpdateUserName(r.Context(), newUser.Username, u.Id)
	if err != nil {
		log.Printf("%s", err.Error())
		http.Error(w, "could not update username", http.StatusInternalServerError)
		return
	}

	h.Hub.UpdateNameChannel(u.Id, newUser.Username)

	w.WriteHeader(http.StatusNoContent)
}

func NewUserHandler(repo interfaces.UserRepository, hub *socket.Hub) *Handler {
	return &Handler{
		Repo: repo,
		Hub:  hub,
	}
}
