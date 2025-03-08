package socket

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	conn     *websocket.Conn
	username string
	hub      *Hub
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

		c.hub.broadcast <- bytes.TrimSpace(message)
	}
}
