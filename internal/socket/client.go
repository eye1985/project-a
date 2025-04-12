package socket

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type MessageJSON struct {
	Message   string    `json:"message"`
	Event     string    `json:"Event"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

type sendMessage struct {
	ClientId   string
	ToChannels []string
	Message    *MessageJSON
}

type client struct {
	id         string
	conn       *websocket.Conn
	email      string
	username   string
	hub        *Hub
	send       chan []byte
	tickerWait time.Duration
	pongWait   time.Duration
	channels   []string
}

const (
	messageTypeMessage = "Message"
	messageTypeJoin    = "join"
	messageTypeQuit    = "quit"
)

type ClientFactory func(conn *websocket.Conn, hub *Hub, username string, channel string) *client

func NewClient(conn *websocket.Conn, hub *Hub, username string, channel string) *client {
	return &client{
		id:         uuid.NewString(),
		conn:       conn,
		username:   username,
		hub:        hub,
		send:       make(chan []byte, 256),
		tickerWait: 10 * time.Second,
		pongWait:   15 * time.Second,
		channels:   []string{channel},
	}
}

func (c *client) read() {
	defer func() {
		quitMsg := &sendMessage{
			ClientId:   c.id,
			ToChannels: c.channels,
			Message: &MessageJSON{
				Message:   c.username + " quit",
				Event:     messageTypeQuit,
				Username:  c.username,
				CreatedAt: time.Now(),
			},
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
				log.Println("Websocket read Message failed: ", err)
			}
			break
		}

		msg := &sendMessage{
			ClientId:   c.id,
			ToChannels: c.channels,
			Message: &MessageJSON{
				Message:   string(message),
				Username:  c.username,
				CreatedAt: time.Now().UTC(),
				Event:     messageTypeMessage,
			},
		}

		message, err = json.Marshal(msg)
		if err != nil {
			log.Printf("Cannot transform json: %s", err)
			break
		}

		c.hub.broadcast <- bytes.TrimSpace(message)
	}
}

func (c *client) write() {
	ticker := time.NewTicker(c.tickerWait)

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte("Closing connection"))
				_ = c.conn.Close()
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				_ = c.conn.Close()
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
