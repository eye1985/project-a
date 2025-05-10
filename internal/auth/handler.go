package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"project-a/internal/contacts"
	"project-a/internal/shared"
)

type Handler struct {
	Repo         Repository
	Service      shared.AuthService
	UserRepo     shared.UserRepository
	UserlistRepo contacts.Repository
}

func createMagicLink(ctx context.Context, email string, h *Handler) (string, error) {
	code, err := h.Service.CreateMagicLink(ctx, email)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		http.Error(w, "code is required", http.StatusBadRequest)
		return
	}

	ml, err := h.Repo.ActivateNonExpiredMagicLink(r.Context(), code)
	if err != nil {
		http.Error(w, "invalid magic link", http.StatusBadRequest)
		return
	}

	username, err := preExtractEmail(ml.Email)
	if err != nil {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	u, err := h.UserRepo.InsertUser(r.Context(), username, ml.Email)
	if err != nil {
		log.Printf("failed to register user: %v", err)
		http.Error(w, "User registration failed", http.StatusInternalServerError)
		return
	}

	err = h.UserlistRepo.CreateContactList(r.Context(), "Contacts", u.Id)
	if err != nil {
		log.Printf("failed to create userlist: %v", err)
		http.Error(w, "Userlist creation failed", http.StatusInternalServerError)
		return
	}

	session, err := h.Service.CreateOrGetSession(r.Context(), u.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoded, err := h.Service.SignCookie(string(shared.SessionCtxKey), []byte(session.SessionID))
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(
		w, &http.Cookie{
			Name:     string(shared.SessionCtxKey),
			Value:    encoded,
			Path:     "/",
			Expires:  session.ExpiresAt,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
	)

	http.Redirect(w, r, shared.HomeRoute, http.StatusSeeOther)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		http.Error(w, "code is required", http.StatusBadRequest)
		return
	}

	ml, err := h.Repo.ActivateNonExpiredMagicLink(r.Context(), code)
	if err != nil {
		http.Error(w, "invalid magic link", http.StatusBadRequest)
		return
	}

	u, err := h.UserRepo.GetUserByEmail(r.Context(), ml.Email)
	if err != nil {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	session, err := h.Service.CreateOrGetSession(r.Context(), u.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoded, err := h.Service.SignCookie(string(shared.SessionCtxKey), []byte(session.SessionID))
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(
		w, &http.Cookie{
			Name:     string(shared.SessionCtxKey),
			Value:    encoded,
			Path:     "/",
			Expires:  session.ExpiresAt,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
	)

	http.Redirect(w, r, shared.HomeRoute, http.StatusSeeOther)
}

// TODO add some security for this
func (h *Handler) CreateMagicLinkCode(w http.ResponseWriter, r *http.Request) {
	var magicLink MagicLink

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&magicLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = mail.ParseAddress(magicLink.Email)
	if err != nil {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	magicCode, err := createMagicLink(r.Context(), magicLink.Email, h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if magicCode == "" {
		http.Error(w, "Could not create magic link", http.StatusInternalServerError)
		return
	}

	// TODO send email
	m := map[string]string{}
	m["magicLinkCode"] = magicCode
	if err := json.NewEncoder(w).Encode(m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Context().Value(shared.SessionCtxKey).([]byte)
	_ = h.Repo.DeleteSession(r.Context(), string(sessionID))

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func NewHandler(
	as shared.AuthService,
	repo Repository,
	ur shared.UserRepository,
	ulr contacts.Repository,
) *Handler {
	return &Handler{
		Repo:         repo,
		Service:      as,
		UserRepo:     ur,
		UserlistRepo: ulr,
	}
}
