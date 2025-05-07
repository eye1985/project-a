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

type sendMessage struct {
	ClientId    int64
	ToClientIds []int64
	Message     *MessageJSON
}

type client struct {
	id         int64
	uuid       uuid.UUID
	conn       *websocket.Conn
	email      string
	username   string
	hub        *Hub
	send       chan []byte
	tickerWait time.Duration
	pongWait   time.Duration
	talkingTo  []int64
	contacts   []int64
}

const (
	messageTypeMessage  = "message"
	messageTypeJoin     = "join"
	messageTypeQuit     = "quit"
	messageTypeIsOnline = "isOnline"
)

type ClientFactory func(
	conn *websocket.Conn,
	hub *Hub,
	id int64,
	username string,
	contacts []int64,
	uuid uuid.UUID,
) *client

func newClient(conn *websocket.Conn, hub *Hub, id int64, username string, contacts []int64, uuid uuid.UUID) *client {
	return &client{
		id:         id,
		uuid:       uuid,
		conn:       conn,
		username:   username,
		hub:        hub,
		send:       make(chan []byte, 256),
		tickerWait: 10 * time.Second,
		pongWait:   15 * time.Second,
		contacts:   contacts,
	}
}

func (c *client) read() {
	defer func() {
		quitMsg := &sendMessage{
			ClientId:    c.id,
			ToClientIds: c.contacts,
			Message: &MessageJSON{
				Message:   c.username + " quit",
				Uuid:      c.uuid,
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

	c.conn.SetPongHandler(
		func(string) error {
			dErr := c.conn.SetReadDeadline(time.Now().Add(c.pongWait))
			if dErr != nil {
				return dErr
			}
			return nil
		},
	)

	c.conn.SetReadLimit(2048)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			var closeErr *websocket.CloseError
			if !errors.As(err, &closeErr) {
				log.Println("Websocket read Message failed: ", err)
			}
			break
		}

		messageIn := &MessageIn{}
		err = json.Unmarshal(message, messageIn)
		if err != nil {
			log.Printf("Cannot unmarshal json: %s", err)
			break
		}

		id, ok := c.hub.uuidToIdMap[messageIn.ToUuid]
		if !ok {
			log.Printf("Cannot find user with uuid: %s", messageIn.ToUuid)
		}

		c.talkingTo = []int64{id, c.id}

		msg := &sendMessage{
			ClientId:    c.id,
			ToClientIds: c.talkingTo,
			Message: &MessageJSON{
				Message:   messageIn.Msg,
				Username:  c.username,
				Uuid:      c.uuid,
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
				log.Printf("Channel closed, closing socket for: %s %d", c.username, c.id)
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte("Closing connection"))
				_ = c.conn.Close()
				return
			}

			// Batch together buffered messages
			messages := [][]byte{msg}
			for i := 0; i < len(c.send); i++ {
				messages = append(messages, <-c.send)
			}

			// Build JSON array
			batched := []byte("[")
			for i, m := range messages {
				batched = append(batched, m...)

				//If this is not the last message, add a comma
				if i < len(messages)-1 {
					batched = append(batched, ',')
				}
			}
			batched = append(batched, ']')

			// Send as a single WebSocket message
			err := c.conn.WriteMessage(websocket.TextMessage, batched)
			if err != nil {
				log.Printf("Writer error, closing socket for: %s %d", c.username, c.id)
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				_ = c.conn.Close()
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
