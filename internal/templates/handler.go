package templates

import (
	"fmt"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
	"project-a/internal/consts"
	"project-a/internal/interfaces"
	"project-a/internal/model"
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
	userRepo     interfaces.UserRepository
	userListRepo interfaces.ContactsRepository
	authService  interfaces.AuthService
	wsUrl        string
	isDev        bool
}

func (h *Handler) RenderRegisterUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(string(consts.SessionCtxKey))
	if err == nil {
		cookieValue, _ := h.authService.VerifyCookie(cookie)
		if h.authService.IsSessionActive(r.Context(), string(cookieValue)) {
			http.Redirect(w, r, consts.HomeRoute, http.StatusSeeOther)
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
		w, &CreateMagicLinkPage{
			PageData: PageData{
				Title: "Register / Login",
			},
			CSRF: csrf.Token(r),
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
	sessionID := r.Context().Value(consts.SessionCtxKey).([]byte)
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

	sessionID := r.Context().Value(consts.SessionCtxKey).([]byte)
	u, _ := h.userRepo.GetUserFromSessionId(r.Context(), string(sessionID))

	contactList, err := h.userListRepo.GetContactLists(r.Context(), u.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	listMap := make(map[*model.List][]*model.Contact)

	for _, ul := range contactList {
		listOfContact, err := h.userListRepo.GetContacts(r.Context(), ul.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		listMap[ul] = listOfContact
	}

	var js string
	var cssList []string

	err = setJsCssPathsFromManifest("chat", &js, &cssList, h.isDev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(
		w, &ChatPage{
			PageData: PageData{
				Title:    "Chat",
				Username: u.Username,
				Uuid:     u.Uuid,
				Js:       js,
				Css:      cssList,
				WsUrl:    h.wsUrl,
			},
			ContactLists: listMap,
		},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RenderContacts(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		fmt.Sprintf("%s/%s", path, baseLayout),
		fmt.Sprintf("%s/%s", path, userContacts),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID := r.Context().Value(consts.SessionCtxKey).([]byte)
	u, _ := h.userRepo.GetUserFromSessionId(r.Context(), string(sessionID))

	contactList, err := h.userListRepo.GetContactLists(r.Context(), u.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	listMap := make(map[*model.List][]*model.Contact)

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

	err = setJsCssPathsFromManifest("contacts", &js, &cssList, h.isDev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(
		w, &ContactPage{
			PageData: PageData{
				Title:    "Contacts",
				Username: u.Username,
				Uuid:     u.Uuid,
				WsUrl:    h.wsUrl,
				Js:       js,
				Css:      cssList,
			},
			ContactLists: listMap,
			Invitations:  invitationTemplates,
		},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewHandler(
	userRepo interfaces.UserRepository,
	userlistRepo interfaces.ContactsRepository,
	authService interfaces.AuthService,
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
