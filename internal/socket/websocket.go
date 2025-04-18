package socket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"project-a/internal/shared"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	channelNameDoesNotExist = "channels does not exist"
)

func ServeWs(hub *Hub, cf ClientFactory, session shared.Session, ur shared.UserRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			log.Printf("origin %s", origin)
			// TODO add origin check here
			return true
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Websocket upgrade failed: ", err)
			return
		}

		cookie, err := r.Cookie(string(shared.SessionCtxKey))
		if err != nil {
			_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "unauthorized"), time.Now().Add(time.Second))
			_ = conn.Close()
			return
		}

		cookieValue, err := session.VerifyCookie(cookie)
		if err != nil {
			_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "unauthorized"), time.Now().Add(time.Second))
			_ = conn.Close()
			return
		}

		if !session.IsSessionActive(string(cookieValue)) {
			_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "unauthorized"), time.Now().Add(time.Second))
			_ = conn.Close()
			return
		}

		channel := r.URL.Query().Get("channels")
		if channel == "" {
			_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, channelNameDoesNotExist), time.Now().Add(time.Second))
			_ = conn.Close()
			return
		}

		u, err := ur.GetUserFromSessionId(string(cookieValue))
		if err != nil {
			_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, channelNameDoesNotExist), time.Now().Add(time.Second))
			_ = conn.Close()
			return
		}

		client := cf(conn, hub, u.Id, u.Username, channel)
		hub.register <- client

		joinMsg := sendMessage{
			ClientId:   client.id,
			ToChannels: client.channels,
			Message: &MessageJSON{
				Message:   client.username + " joined",
				Event:     messageTypeJoin,
				Username:  client.username,
				CreatedAt: time.Now(),
			},
		}

		joinedMsgByte, err := json.Marshal(joinMsg)
		if err != nil {
			log.Println("Websocket marshal failed: ", err)
		}

		hub.broadcast <- joinedMsgByte

		go client.read()
		go client.write()
	}
}
