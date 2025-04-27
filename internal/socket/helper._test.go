package socket

import (
	"errors"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"project-a/internal/shared"
	"time"
)

type HubAndTestServer struct {
	server *httptest.Server
	hub    *Hub
}

type MockSession struct{}

func (m *MockSession) VerifyCookie(cookie *http.Cookie) ([]byte, error) {
	return []byte("1234567890"), nil
}

func (m *MockSession) IsSessionActive(sessionId string) bool {
	return true
}

func (m *MockSession) SignCookie(cookieName string, value []byte) (string, error) {
	return "1234567890", nil
}

func (h *HubAndTestServer) wsUrl(username string, channel string) string {
	return "ws" + h.server.URL[len("http"):] + "?username=" + username + "&channels=" + channel
}

type MockUserRepo struct {
	user *shared.User
}

func (r *MockUserRepo) GetUser(email string) (*shared.User, error) {
	return &shared.User{}, nil
}

func (r *MockUserRepo) GetUserFromSessionId(sessionId string) (*shared.User, error) {
	return &shared.User{}, nil
}

func (r *MockUserRepo) GetUsers() ([]*shared.User, error) {
	return []*shared.User{}, nil
}

func (r *MockUserRepo) InsertUser(username string, email string) (*shared.User, error) {
	return &shared.User{}, nil
}

func (r *MockUserRepo) DeleteUser(email string) error {
	return nil
}

func (r *MockUserRepo) UpdateUserName(username string, userId int64) error {
	return nil
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

	//mockSession := &MockSession{}
	//mockUserRepo := &MockUserRepo{
	//	user: &shared.User{},
	//}

	//wsHandler := ServeWs(hub, cf, mockSession, mockUserRepo)
	return HubAndTestServer{
		//server: httptest.NewServer(http.HandlerFunc(wsHandler)),
		hub: hub,
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
