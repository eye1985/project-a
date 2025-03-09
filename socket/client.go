package socket

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	conn     *websocket.Conn
	username string
	hub      *Hub
	send     chan []byte
}

type MessageJSON struct {
	Message   string    `json:"message"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

func (c *Client) read() {
	defer func() {
		log.Println("Disconnect from websocket")
		c.hub.unregister <- c
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Websocket read message failed: ", err)
			break
		}

		msg := &MessageJSON{
			Message:   string(message),
			Username:  c.username,
			CreatedAt: time.Now().UTC(),
		}

		message, err = json.Marshal(msg)
		if err != nil {
			log.Printf("Cannot transform json: %s", err)
			break
		}

		c.hub.broadcast <- bytes.TrimSpace(message)
	}
}
