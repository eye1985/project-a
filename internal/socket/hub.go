package socket

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
)

type ClientNameUpdateRequest struct {
	ClientId int64
	Username string
}

type Hub struct {
	broadcast   chan []byte
	register    chan *client
	unregister  chan *client
	clients     map[int64]*client
	uuidToIdMap map[uuid.UUID]int64
	updateName  chan ClientNameUpdateRequest
}

func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *client),
		unregister:  make(chan *client),
		clients:     make(map[int64]*client),
		uuidToIdMap: make(map[uuid.UUID]int64),
		updateName:  make(chan ClientNameUpdateRequest),
	}
}

func (h *Hub) UpdateNameChannel(clientId int64, username string) {
	h.updateName <- ClientNameUpdateRequest{
		ClientId: clientId,
		Username: username,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case userClient := <-h.register:
			if existing, ok := h.clients[userClient.id]; ok {
				_ = existing.conn.WriteMessage(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing"),
				)
				_ = h.clients[userClient.id].conn.Close()
			}

			h.clients[userClient.id] = userClient
			h.uuidToIdMap[userClient.uuid] = userClient.id
		case userClient := <-h.unregister:
			_ = userClient.conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing"),
			)
			_ = userClient.conn.Close()
			delete(h.clients, userClient.id)
		case message := <-h.broadcast:
			var msg sendMessage
			err := json.Unmarshal(message, &msg)
			if err != nil {
				log.Printf("Error unmarshalling Message: %s", err)
				continue
			}
			for _, clientId := range msg.ToClientIds {

				log.Printf("Sending message to %d", clientId)

				jsonMsg, err := json.Marshal(msg.Message)
				if err != nil {
					log.Printf("Error marshalling Message: %s", err)
					continue
				}

				ignoreSelf := msg.Message.Event == messageTypeJoin || msg.Message.Event == messageTypeQuit
				if ignoreSelf && clientId == msg.ClientId {
					log.Printf("Ignoring self message")
					continue
				}

				if target, ok := h.clients[clientId]; ok {
					target.send <- jsonMsg
					log.Printf("Message sent to %d", clientId)
				}
			}
		case clientReq := <-h.updateName:
			if c, ok := h.clients[clientReq.ClientId]; ok {
				c.username = clientReq.Username
			}
		}
	}
}
