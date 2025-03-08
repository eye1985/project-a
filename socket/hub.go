package socket

import (
	"github.com/gorilla/websocket"
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
			for username, clients := range h.clients {
				for _, client := range clients {

					// TODO figure out batch sending, using NextWriter and chan for sending. With channels you can check queue
					if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
						_ = client.conn.Close()
						h.clients[username] = slices.DeleteFunc(h.clients[username], func(c *Client) bool {
							return c == client
						})

						log.Printf("Cant write %s %s", username, err)
					}
				}

				if len(h.clients[username]) == 0 {
					delete(h.clients, username)
				}
			}

			log.Printf("Client broadcast: %s", string(message))
		}
	}
}
