package templates

import (
	"fmt"
	"html/template"
	"net/http"
	"project-a/internal/contacts"
	"project-a/internal/shared"
)

const (
	path         = "web/templates"
	manifestPath = "web/assets/dist/manifest.json"
	prodPath     = "assets/dist/"
	devJsPath    = "assets/js/"
	devCssPath   = "assets/css/"
	baseLayout   = "base-layout.gohtml"
	register     = "register.gohtml"
	profile      = "profile.gohtml"
	userContacts = "contacts.gohtml"
	chat         = "chat.gohtml"
)

type Handler struct {
	userRepo     shared.UserRepository
	userListRepo contacts.Repository
	authService  shared.AuthService
	wsUrl        string
	isDev        bool
}

func (h *Handler) RenderRegisterUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(string(shared.SessionCtxKey))
	if err == nil {
		cookieValue, _ := h.authService.VerifyCookie(cookie)
		if h.authService.IsSessionActive(r.Context(), string(cookieValue)) {
			http.Redirect(w, r, shared.HomeRoute, http.StatusSeeOther)
			return
		}
	}

	tmpl, err := template.ParseFiles(
		fmt.Sprintf("%s/%s", path, register),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(
		w, &PageData{
			Title: "Profile",
		},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RenderProfile(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		fmt.Sprintf("%s/%s", path, baseLayout),
		fmt.Sprintf("%s/%s", path, profile),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sessionID := r.Context().Value(shared.SessionCtxKey).([]byte)
	u, err := h.userRepo.GetUserFromSessionId(r.Context(), string(sessionID))
	if err != nil {
		http.Error(w, "no session", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(
		w, &PageData{
			Username: u.Username,
			Title:    "Profile",
		},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RenderChat(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		fmt.Sprintf("%s/%s", path, baseLayout),
		fmt.Sprintf("%s/%s", path, chat),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID := r.Context().Value(shared.SessionCtxKey).([]byte)
	u, _ := h.userRepo.GetUserFromSessionId(r.Context(), string(sessionID))

	contactList, err := h.userListRepo.GetContactLists(r.Context(), u.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	listMap := make(map[*contacts.List][]*contacts.Contact)

	for _, ul := range contactList {
		listOfContact, err := h.userListRepo.GetContacts(r.Context(), ul.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		listMap[ul] = listOfContact
	}

	invitations, err := h.userListRepo.GetInvitations(r.Context(), u.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var invitationTemplates []*InvitationTemplate
	for _, inv := range invitations {
		var invitationTemplate InvitationTemplate
		invitationTemplate.InviteUuid = inv.Uuid

		if u.Id == inv.InviteeId {
			invitationTemplate.IsInviter = false
			invitationTemplate.Email = inv.InviterEmail

		} else if u.Id == inv.InviterId {
			invitationTemplate.IsInviter = true
			invitationTemplate.Email = inv.InviteeEmail
		}

		invitationTemplates = append(invitationTemplates, &invitationTemplate)
	}

	var js string
	var cssList []string

	err = setJsCssPathsFromManifest(&js, &cssList, h.isDev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(
		w, &ChatPage{
			Title:        "Chat",
			Username:     u.Username,
			Uuid:         u.Uuid,
			ContactLists: listMap,
			Invitations:  invitationTemplates,
			WsUrl:        h.wsUrl,
			Js:           js,
			Css:          cssList,
		},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewHandler(
	userRepo shared.UserRepository,
	userlistRepo contacts.Repository,
	authService shared.AuthService,
	wsUrl string,
	isDev bool,
) *Handler {
	return &Handler{
		userRepo:     userRepo,
		userListRepo: userlistRepo,
		authService:  authService,
		wsUrl:        wsUrl,
		isDev:        isDev,
	}
}
