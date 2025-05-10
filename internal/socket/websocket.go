package socket

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"project-a/internal/contacts"
	"project-a/internal/shared"
	"time"
)

const (
	sessionDoesNotExist = "session does not exist"
	contactError        = "contact error"
)

func ServeWs(
	hub *Hub,
	cf ClientFactory,
	as shared.AuthService,
	ur shared.UserRepository,
	cr contacts.Repository,
	origin string,
) func(
	http.ResponseWriter,
	*http.Request,
) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			rOrigin := r.Header.Get("Origin")
			log.Printf("rOrigin: %s", rOrigin)
			if rOrigin != origin {
				return false
			}
			return true
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Websocket upgrade failed: ", err)
			return
		}

		cookie, err := r.Cookie(string(shared.SessionCtxKey))
		if err != nil {
			_ = conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "unauthorized"),
				time.Now().Add(time.Second),
			)
			_ = conn.Close()
			return
		}

		cookieValue, err := as.VerifyCookie(cookie)
		if err != nil {
			_ = conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "unauthorized"),
				time.Now().Add(time.Second),
			)
			_ = conn.Close()
			return
		}

		if !as.IsSessionActive(r.Context(), string(cookieValue)) {
			_ = conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "unauthorized"),
				time.Now().Add(time.Second),
			)
			_ = conn.Close()
			return
		}

		u, err := ur.GetUserFromSessionId(r.Context(), string(cookieValue))
		if err != nil {
			_ = conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, sessionDoesNotExist),
				time.Now().Add(time.Second),
			)
			_ = conn.Close()
			return
		}

		listOfContact, err := cr.GetContacts(r.Context(), u.Id)
		if err != nil {
			_ = conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, contactError),
				time.Now().Add(time.Second),
			)
			_ = conn.Close()
			return
		}

		var contactIds []int64
		contactsOnline := []uuid.UUID{}
		for _, c := range listOfContact {
			contactIds = append(contactIds, c.UserId)
			_, ok := hub.clients[c.UserId]
			if ok {
				contactsOnline = append(contactsOnline, c.UserUuid)
			}
		}

		client := cf(conn, hub, u.Id, u.Username, contactIds, u.Uuid)
		hub.register <- client

		joinMsg := sendMessage{
			ClientId:    client.id,
			ToClientIds: contactIds,
			Message: &MessageJSON{
				FromUuid:  client.uuid,
				Message:   client.username + " joined",
				Event:     messageTypeJoin,
				Username:  client.username,
				CreatedAt: time.Now(),
			},
		}

		isOnlineList, err := json.Marshal(contactsOnline)
		if err != nil {
			log.Println("Websocket marshal failed: ", err)
		}

		listOfUserContactsOnlineMsg := sendMessage{
			ClientId:    client.id,
			ToClientIds: []int64{client.id},
			Message: &MessageJSON{
				FromUuid:  client.uuid,
				Message:   string(isOnlineList),
				Event:     messageTypeIsOnline,
				Username:  client.username,
				CreatedAt: time.Now(),
			},
		}

		joinedMsgByte, err := json.Marshal(joinMsg)
		if err != nil {
			log.Println("Websocket marshal failed: ", err)
		}

		hub.broadcast <- joinedMsgByte

		listOfUserContactsOnlineMsgByte, err := json.Marshal(listOfUserContactsOnlineMsg)
		if err != nil {
			log.Println("Websocket marshal failed: ", err)
		}
		hub.broadcast <- listOfUserContactsOnlineMsgByte

		go client.read()
		go client.write()
	}
}
