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

const (
	UsernameDoesNotExist = "username does not exist"
)

func ServeWs(hub *Hub, cf ClientFactory) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Websocket upgrade failed: ", err)
			return
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, UsernameDoesNotExist), time.Now().Add(time.Second))
			_ = conn.Close()
			return
		}

		log.Println("Websocket connected: ", username)
		client := cf(conn, hub, username)
		hub.register <- client
		go client.read()
		go client.write()
	}
}
