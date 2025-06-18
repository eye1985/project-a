package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"project-a/internal/consts"
	"project-a/internal/email"
	"project-a/internal/interfaces"
	"project-a/internal/util"
	"strings"
)

type Handler struct {
	Repo          Repository
	Service       interfaces.AuthService
	UserRepo      interfaces.UserRepository
	UserlistRepo  interfaces.ContactsRepository
	EmailRepo     email.Repository
	MailSendToken string
	Origin        string
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

	encoded, err := h.Service.SignCookie(string(consts.SessionCtxKey), []byte(session.SessionID))
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(
		w, &http.Cookie{
			Name:     string(consts.SessionCtxKey),
			Value:    encoded,
			Path:     "/",
			Expires:  session.ExpiresAt,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
	)

	http.Redirect(w, r, consts.HomeRoute, http.StatusSeeOther)
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

	encoded, err := h.Service.SignCookie(string(consts.SessionCtxKey), []byte(session.SessionID))
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(
		w, &http.Cookie{
			Name:     string(consts.SessionCtxKey),
			Value:    encoded,
			Path:     "/",
			Expires:  session.ExpiresAt,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
	)

	http.Redirect(w, r, consts.HomeRoute, http.StatusSeeOther)
}

// TODO add some security for this
func (h *Handler) CreateMagicLinkCode(w http.ResponseWriter, r *http.Request) {
	var magicLink MagicLink

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&magicLink)
	if err != nil {
		http.Error(w, "Missing body", http.StatusBadRequest)
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
	ip, err := util.ReadUserIP(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u, _ := h.UserRepo.GetUserByEmail(r.Context(), magicLink.Email)
	isSignUp := u == nil
	err = h.EmailRepo.AddSentEmail(r.Context(), magicLink.Email, ip, isSignUp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var name string
	if u != nil {
		name = u.Username
	} else {
		name = strings.Split(magicLink.Email, "@")[0]
	}

	err = email.SendEmailFromMailSend(
		h.MailSendToken,
		&email.SendEmailFromMailSendArgs{
			Email:    magicLink.Email,
			Name:     name,
			IsSignUp: isSignUp,
			Code:     magicCode,
			Origin:   h.Origin,
		},
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// For debugging
	//m := map[string]string{}
	//m["magicLinkCode"] = magicCode
	//if err := json.NewEncoder(w).Encode(m); err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Context().Value(consts.SessionCtxKey).([]byte)
	_ = h.Repo.DeleteSession(r.Context(), string(sessionID))

	http.SetCookie(
		w, &http.Cookie{
			Name:     string(consts.SessionCtxKey),
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
	)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type NewHandlerArgs struct {
	AuthService   interfaces.AuthService
	Repo          Repository
	UserRepo      interfaces.UserRepository
	ContactsRepo  interfaces.ContactsRepository
	EmailRepo     email.Repository
	MailSendToken string
	Origin        string
}

func NewHandler(args *NewHandlerArgs) *Handler {
	return &Handler{
		Repo:          args.Repo,
		Service:       args.AuthService,
		UserRepo:      args.UserRepo,
		UserlistRepo:  args.ContactsRepo,
		EmailRepo:     args.EmailRepo,
		MailSendToken: args.MailSendToken,
		Origin:        args.Origin,
	}
}
