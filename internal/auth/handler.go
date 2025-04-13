package auth

import (
	"net/http"
	"net/mail"
	"project-a/internal/user"
)

type Handler struct {
	Service
	UserService user.Service
}

func registerIfUserNotExists(email string, username string, h *Handler) (*user.User, error) {
	u, err := h.UserService.GetUser(email)
	if err != nil {
		newUser, err := h.UserService.RegisterUser(&user.User{
			Username: username,
			Email:    email,
		})
		if err != nil {
			return nil, err
		}

		return newUser, nil
	}
	return u, nil
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

	username, err := preExtractEmail(email)
	if err != nil {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	newUser, err := registerIfUserNotExists(email, username, h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := h.Service.CreateOrGetSession(newUser.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookieName := "sid"
	encoded, err := h.Service.SignCookie(cookieName, []byte(session.SessionID))
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    encoded,
		Path:     "/",
		Expires:  session.ExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

func NewAuthHandler(as Service, us user.Service) *Handler {
	return &Handler{
		Service:     as,
		UserService: us,
	}
}
