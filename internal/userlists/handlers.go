package userlists

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	Repo UserListRepository
}

func (h *Handler) CreateUserList(w http.ResponseWriter, r *http.Request) {
	var body CreateUserListBody
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repo.CreateUserList(r.Context(), body.Name, body.UserId)
	if err != nil {
		http.Error(w, "Could not create userlist", http.StatusInternalServerError)
		return
	}
}

func NewHandler(repo UserListRepository) *Handler {
	return &Handler{repo}
}
