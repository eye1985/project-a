package user

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	Repo Repository
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

	user := &User{}
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

func NewUserHandler(repo Repository) *Handler {
	return &Handler{
		Repo: repo,
	}
}
