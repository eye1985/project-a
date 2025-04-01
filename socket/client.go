package socket

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	conn       *websocket.Conn
	username   string
	hub        *Hub
	send       chan []byte
	tickerWait time.Duration
	pongWait   time.Duration
	channel    []string // TODO implement later
}

const (
	messageTypeMessage = "message"
	messageTypeJoin    = "join"
	messageTypeQuit    = "quit"
)

type ClientFactory func(conn *websocket.Conn, hub *Hub, username string) *Client

func NewClient(conn *websocket.Conn, hub *Hub, username string) *Client {
	return &Client{
		conn:       conn,
		username:   username,
		hub:        hub,
		send:       make(chan []byte, 256),
		tickerWait: 10 * time.Second,
		pongWait:   15 * time.Second,
		channel:    []string{},
	}
}

type MessageJSON struct {
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

func (c *Client) read() {
	defer func() {
		quitMsg := MessageJSON{
			Message:   c.username + " quit",
			Type:      messageTypeQuit,
			Username:  c.username,
			CreatedAt: time.Now(),
		}

		msg, err := json.Marshal(quitMsg)
		if err != nil {
			c.hub.unregister <- c
		} else {
			c.hub.unregister <- c
			c.hub.broadcast <- msg
		}
	}()

	dErr := c.conn.SetReadDeadline(time.Now().Add(c.pongWait))
	if dErr != nil {
		return
	}

	c.conn.SetPongHandler(func(string) error {
		dErr := c.conn.SetReadDeadline(time.Now().Add(c.pongWait))
		if dErr != nil {
			return dErr
		}
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			var closeErr *websocket.CloseError
			if !errors.As(err, &closeErr) {
				log.Println("Websocket read message failed: ", err)
			}
			break
		}

		msg := &MessageJSON{
			Message:   string(message),
			Username:  c.username,
			CreatedAt: time.Now().UTC(),
			Type:      messageTypeMessage,
		}

		message, err = json.Marshal(msg)
		if err != nil {
			log.Printf("Cannot transform json: %s", err)
			break
		}

		c.hub.broadcast <- bytes.TrimSpace(message)
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(c.tickerWait)

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte("Closing connection"))
				c.conn.Close()
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.conn.Close()
				return
			}
			newline := []byte{'\n'}
			w.Write(msg)
			w.Write(newline)
			for i := 0; i < len(c.send); i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			err = w.Close()
			if err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
