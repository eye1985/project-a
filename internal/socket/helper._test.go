package socket

import (
	"errors"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

type HubAndTestServer struct {
	server *httptest.Server
	hub    *Hub
}

func (h *HubAndTestServer) wsUrl(username string, channel string) string {
	return "ws" + h.server.URL[len("http"):] + "?username=" + username + "&channels=" + channel
}

func newTestClient(conn *websocket.Conn, hub *Hub, username string, channel string) *client {
	return &client{
		conn:       conn,
		username:   username,
		hub:        hub,
		send:       make(chan []byte, 256),
		tickerWait: 10 * time.Millisecond,
		pongWait:   15 * time.Millisecond,
		channels:   []string{channel},
	}
}

func testServerHub(silence bool, cf ClientFactory) HubAndTestServer {
	if silence {
		log.SetOutput(io.Discard)
	}
	hub := NewHub()
	go hub.Run()

	wsHandler := ServeWs(hub, cf)
	return HubAndTestServer{
		server: httptest.NewServer(http.HandlerFunc(wsHandler)),
		hub:    hub,
	}
}

type pollingForClientInChannelsArgs struct {
	sec     time.Duration
	hub     *Hub
	channel string
	cap     int
}

func pollingForClientInChannels(props pollingForClientInChannelsArgs) error {
	timeout := time.After(props.sec * time.Second)
	tick := time.Tick(10 * time.Millisecond)

	for {
		select {
		case <-timeout:
			return errors.New("timeout waiting for client")
		case <-tick:
			if len(props.hub.channels[props.channel]) == props.cap {
				return nil
			}
		}
	}
}
