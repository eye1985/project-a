package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type webSocket struct {
	mutex   sync.Mutex
	clients map[string]*wsConn
}

func (ws *webSocket) usernameExists(username string) bool {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	_, key := ws.clients[username]
	return key
}

func (ws *webSocket) addClient(username string, wsConn *wsConn) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	log.Println("Adding client:", wsConn.client.RemoteAddr().String())
	ws.clients[username] = wsConn
}

func (ws *webSocket) removeClient(username string, wsConn *wsConn) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	log.Printf("Removing client %s", wsConn.client.RemoteAddr().String())

	conn, exist := ws.clients[username]
	if !exist {
		return fmt.Errorf("user %s does not exist", username)
	}

	if conn != wsConn {
		return fmt.Errorf("user %s has a different active connection", username)
	}

	err := conn.client.Close()
	if err != nil {
		return err
	}

	delete(ws.clients, username)
	return nil
}

func (ws *webSocket) broadcast(msg []byte, messageType int) {
	ws.mutex.Lock()
	clientsCopy := make(map[string]*wsConn, len(ws.clients))
	for username, conn := range ws.clients {
		clientsCopy[username] = conn
	}
	ws.mutex.Unlock()

	var usernames []string
	for username, conn := range clientsCopy {
		if err := conn.client.WriteMessage(messageType, msg); err != nil {
			log.Println("Cannot send message to client, closing client ", err)
			_ = conn.client.Close()
			usernames = append(usernames, username)
		}
	}

	ws.mutex.Lock()
	for _, username := range usernames {
		delete(ws.clients, username)
	}
	ws.mutex.Unlock()
}

type wsConn struct {
	client *websocket.Conn
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ws = &webSocket{
	mutex:   sync.Mutex{},
	clients: make(map[string]*wsConn),
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
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

	if ws.usernameExists(username) {
		_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "username already exists"), time.Now().Add(time.Second))
		_ = conn.Close()
		return
	}

	newWsConn := &wsConn{
		client: conn,
	}

	ws.addClient(username, newWsConn)

	defer func() {
		err = ws.removeClient(username, newWsConn)
		if err != nil {
			log.Println("Cannot close client ", err)
		}
	}()

	for {
		messageType, message, connErr := conn.ReadMessage()
		if connErr != nil {
			log.Println("Websocket read message failed: ", err)
		}

		ws.broadcast([]byte(username+": "+string(message)), messageType)
	}
}
