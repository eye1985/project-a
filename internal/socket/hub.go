package socket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"slices"
)

type Hub struct {
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
	channels   map[string][]*client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
		channels:   make(map[string][]*client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case userClient := <-h.register:
			for _, ch := range userClient.channels {
				h.channels[ch] = append(h.channels[ch], userClient)
				log.Printf("client registered %s and join %s channel", userClient.username, ch)
			}
		case userClient := <-h.unregister:
			_ = userClient.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing"))
			_ = userClient.conn.Close()

			for _, ch := range userClient.channels {
				h.channels[ch] = slices.DeleteFunc(h.channels[ch], func(c *client) bool {
					return c.id == userClient.id
				})
				log.Printf("Channel %s length: %d", ch, len(h.channels[ch]))
			}

			log.Printf("client unregistered %s", userClient.username)
		case message := <-h.broadcast:
			var msg sendMessage
			err := json.Unmarshal(message, &msg)
			if err != nil {
				log.Printf("Error unmarshalling Message: %s", err)
				continue
			}

			for _, ch := range msg.ToChannels {
				for _, c := range h.channels[ch] {
					jsonMsg, err := json.Marshal(msg.Message)
					if err != nil {
						log.Printf("Error marshalling Message: %s", err)
						continue
					}
					ignoreSelf := msg.Message.Event == messageTypeJoin || msg.Message.Event == messageTypeQuit
					if ignoreSelf && c.id == msg.ClientId {
						continue
					}

					c.send <- jsonMsg
					log.Printf("client broadcast: %s", string(jsonMsg))
				}
			}
		}
	}
}
