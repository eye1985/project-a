package socket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs() func(http.ResponseWriter, *http.Request) {
	hub := NewHub()
	go hub.Run()

	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Websocket upgrade failed: ", err)
			return
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "username does not exist"), time.Now().Add(time.Second))
			_ = conn.Close()
			return
		}

		// TODO finish usage of channels
		client := &Client{conn: conn, username: username, hub: hub, send: make(chan []byte, 256)}
		hub.register <- client
		go client.read()
	}
}
