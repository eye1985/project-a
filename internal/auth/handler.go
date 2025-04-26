package auth

import (
	"log"
	"net/http"
	"net/mail"
	"project-a/internal/shared"
)

type Handler struct {
	Repo     Repository
	Service  Service
	UserRepo shared.UserRepository
}

func createMagicLinkIfNotExist(email string, h *Handler) (*shared.User, string, error) {
	u, err := h.UserRepo.GetUser(email)
	if err != nil {
		code, err := h.Service.CreateMagicLink(email)
		if err != nil {
			return nil, "", err
		}

		return nil, code, nil
	}
	return u, "", nil
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		http.Error(w, "code is required", http.StatusBadRequest)
		return
	}

	ml, err := h.Repo.ActivateNonExpiredMagicLink(code)
	if err != nil {
		http.Error(w, "invalid magic link", http.StatusBadRequest)
		return
	}

	username, err := preExtractEmail(ml.Email)
	if err != nil {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	u, err := h.UserRepo.InsertUser(username, ml.Email)
	if err != nil {
		log.Printf("failed to register user: %v", err)
		http.Error(w, "User registration failed", http.StatusInternalServerError)
		return
	}

	session, err := h.Service.CreateOrGetSession(u.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoded, err := h.Service.SignCookie(string(shared.SessionCtxKey), []byte(session.SessionID))
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     string(shared.SessionCtxKey),
		Value:    encoded,
		Path:     "/",
		Expires:  session.ExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, shared.HomeRoute, http.StatusSeeOther)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	_, err = mail.ParseAddress(email)
	if err != nil {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	existingUser, magicCode, err := createMagicLinkIfNotExist(email, h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if magicCode != "" {
		// TODO Send email here
		http.Redirect(w, r, "/?magicLinkCode="+magicCode, http.StatusSeeOther)
		return
	}

	session, err := h.Service.CreateOrGetSession(existingUser.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO create link for login
	encoded, err := h.Service.SignCookie(string(shared.SessionCtxKey), []byte(session.SessionID))
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     string(shared.SessionCtxKey),
		Value:    encoded,
		Path:     "/",
		Expires:  session.ExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, shared.HomeRoute, http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.Context().Value(shared.SessionCtxKey).([]byte)
	_ = h.Repo.DeleteSession(string(sessionID))

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func NewAuthHandler(as Service, repo Repository, ur shared.UserRepository) *Handler {
	return &Handler{
		Repo:     repo,
		Service:  as,
		UserRepo: ur,
	}
}
