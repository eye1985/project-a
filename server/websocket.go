package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type webSocket struct {
	mutex   sync.Mutex
	clients map[*websocket.Conn]bool
}

func (ws *webSocket) addClient(c *websocket.Conn) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	log.Println("Adding client:", c.RemoteAddr().String())
	ws.clients[c] = true
}

func (ws *webSocket) removeClient(c *websocket.Conn) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	log.Printf("Removing client %s", c.RemoteAddr().String())
	delete(ws.clients, c)
}

func (ws *webSocket) broadcast(msg []byte, messageType int) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	for c := range ws.clients {
		if err := c.WriteMessage(messageType, msg); err != nil {
			log.Println("Cannot send message to client, closing client ", err)
			err := c.Close()
			if err != nil {
				log.Println("Cannot close client ", err)
			}
			ws.removeClient(c)
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ws = &webSocket{
	mutex:   sync.Mutex{},
	clients: make(map[*websocket.Conn]bool),
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade failed: ", err)
		return
	}

	ws.addClient(conn)

	defer func() {
		ws.removeClient(conn)
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

		ws.broadcast(message, messageType)
	}
}
