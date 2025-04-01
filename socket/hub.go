package socket

import (
	"encoding/json"
	"log"
	"slices"
)

type Hub struct {
	clients    map[string][]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string][]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.username] = append(h.clients[client.username], client)
			log.Printf("Client registered %s", client.username)
		case client := <-h.unregister:
			for _, c := range h.clients[client.username] {
				if c == client {
					_ = c.conn.Close()
					h.clients[client.username] = slices.DeleteFunc(h.clients[client.username], func(cSlice *Client) bool {
						return cSlice == c
					})
				}
			}

			if len(h.clients[client.username]) == 0 {
				delete(h.clients, client.username)
			}

			log.Printf("Client unregistered %s", client.username)
		case message := <-h.broadcast:
			var messageJSON MessageJSON
			err := json.Unmarshal(message, &messageJSON)
			if err != nil {
				log.Printf("Error unmarshalling message: %s", err)
				continue
			}

			ignoreSelf := messageJSON.Type == messageTypeJoin || messageJSON.Type == messageTypeQuit

			for username, clients := range h.clients {
				for _, client := range clients {
					if ignoreSelf && client.username == messageJSON.Username {
						continue
					}
					client.send <- message
				}
				if len(h.clients[username]) == 0 {
					delete(h.clients, username)
				}
			}

			log.Printf("Client broadcast: %s", string(message))
		}
	}
}
