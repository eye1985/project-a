package contacts

import (
	"encoding/json"
	"log"
	"net/http"
	"project-a/internal/shared"
)

type Handler struct {
	Repo     Repository
	UserRepo shared.UserRepository
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

	err = h.Repo.CreateContactList(r.Context(), body.Name, body.UserId)
	if err != nil {
		http.Error(w, "Could not create userlist", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateContact(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(shared.SessionCtxKey).([]byte)
	inviter, err := h.UserRepo.GetUserFromSessionId(r.Context(), string(session))
	if err != nil {
		http.Error(w, "no user", http.StatusInternalServerError)
		return
	}

	var body CreateContactBody
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Email == "" || body.ContactListId == 0 {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	invitee, err := h.UserRepo.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		http.Error(w, "Could not find invitee", http.StatusInternalServerError)
		return
	}

	_, err = h.Repo.CreateContact(r.Context(), invitee.Id, inviter.Id, body.ContactListId, invitee.Username)
	if err != nil {
		log.Printf("create contact error: %s", err.Error())
		http.Error(w, "Could not create contact", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(shared.SessionCtxKey).([]byte)
	invitee, err := h.UserRepo.GetUserFromSessionId(r.Context(), string(session))
	if err != nil {
		http.Error(w, "no user", http.StatusInternalServerError)
		return
	}

	var body AcceptInviteBody
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repo.UpdateContact(r.Context(), true, body.Uuid, invitee.Id)
	if err != nil {
		log.Printf("accept invite error: %s", err.Error())
		http.Error(w, "Could not accept invite", http.StatusInternalServerError)
	}
}

func NewHandler(repo Repository, userRepo shared.UserRepository) *Handler {
	return &Handler{
		Repo:     repo,
		UserRepo: userRepo,
	}
}
