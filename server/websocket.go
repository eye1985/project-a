package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade failed: ", err)
		return
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("Websocket close failed: ", err)
		}
		log.Println("Websocket closed", r.Host)
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Websocket read message failed: ", err)
			break
		}

		log.Println("Websocket read message: ", string(message))
		log.Println("Websocket message type: ", messageType)

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Websocket write message failed: ", err)
			break
		}
	}
}
