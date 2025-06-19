package contacts

import (
	"encoding/json"
	"log"
	"net/http"
	"project-a/internal/consts"
	"project-a/internal/interfaces"
)

type Handler struct {
	Repo     interfaces.ContactsRepository
	UserRepo interfaces.UserRepository
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

	_, err = h.Repo.CreateContactList(r.Context(), body.Name, body.UserId)
	if err != nil {
		http.Error(w, "Could not create userlist", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateInvitation(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(consts.SessionCtxKey).([]byte)
	inviter, err := h.UserRepo.GetUserFromSessionId(r.Context(), string(session))
	if err != nil {
		http.Error(w, "No user", http.StatusInternalServerError)
		return
	}

	var body CreateInvitationBody
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Email == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if inviter.Email == body.Email {
		http.Error(w, "Cannot invite yourself", http.StatusBadRequest)
		return
	}

	invitee, err := h.UserRepo.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		http.Error(w, "Could not find invitee", http.StatusInternalServerError)
		return
	}

	contacts, err := h.Repo.GetContacts(r.Context(), inviter.Id)
	if err != nil {
		log.Printf("get contacts error: %s", err.Error())
		http.Error(w, "Could not get contacts", http.StatusInternalServerError)
		return
	}

	// Check if user already in your contact list
	for _, c := range contacts {
		if c.UserId == invitee.Id {
			http.Error(w, "Invitee already in your contacts", http.StatusBadRequest)
			return
		}
	}

	err = h.Repo.InviteUser(r.Context(), inviter.Id, invitee.Id)
	if err != nil {
		log.Printf("create contact error in CreateInvitation: %s", err.Error())
		http.Error(w, "Invitation already exists", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(consts.SessionCtxKey).([]byte)
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

	acceptedInvite, err := h.Repo.AcceptInvite(r.Context(), body.Uuid, invitee.Id)
	if err != nil {
		log.Printf("accept invite error: %s", err.Error())
		http.Error(w, "Could not accept invite", http.StatusInternalServerError)
		return
	}

	if acceptedInvite == nil {
		http.Error(w, "Could not accept invite", http.StatusInternalServerError)
		return
	}

	cl1, err := h.Repo.GetContactLists(r.Context(), acceptedInvite.InviterId)
	if err != nil {
		log.Printf("get contact lists error: %s", err.Error())
		http.Error(w, "Could not get contact lists", http.StatusInternalServerError)
		return
	}
	cl2, err := h.Repo.GetContactLists(r.Context(), acceptedInvite.InviteeId)
	if err != nil {
		log.Printf("get contact lists error: %s", err.Error())
		http.Error(w, "Could not get contact lists", http.StatusInternalServerError)
		return
	}

	insertedContact, err := h.Repo.CreateContact(
		r.Context(),
		acceptedInvite.InviterId,
		acceptedInvite.InviteeId,
	)
	if err != nil {
		log.Printf("create contact error in accept invite: %s", err.Error())
		http.Error(w, "Could not create contact", http.StatusInternalServerError)
		return
	}

	err = h.Repo.CreateContactLink(r.Context(), insertedContact.Id, cl1[0].Id)
	if err != nil {
		log.Printf("create contact link error: %s", err.Error())
		http.Error(w, "Could not create contact link", http.StatusInternalServerError)
		return
	}
	err = h.Repo.CreateContactLink(r.Context(), insertedContact.Id, cl2[0].Id)
	if err != nil {
		log.Printf("create contact link error 2: %s", err.Error())
		http.Error(w, "Could not create contact link 2", http.StatusInternalServerError)
		return
	}
}

func NewHandler(
	repo interfaces.ContactsRepository,
	userRepo interfaces.UserRepository,
) *Handler {
	return &Handler{
		Repo:     repo,
		UserRepo: userRepo,
	}
}
