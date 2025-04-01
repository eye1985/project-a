package socket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type HubAndTestServer struct {
	server *httptest.Server
	hub    *Hub
}

func newTestClient(conn *websocket.Conn, hub *Hub, username string) *Client {
	return &Client{
		conn:       conn,
		username:   username,
		hub:        hub,
		send:       make(chan []byte, 256),
		tickerWait: 10 * time.Millisecond,
		pongWait:   15 * time.Millisecond,
		channel:    []string{},
	}
}

func testStruct(silence bool, cf ClientFactory) HubAndTestServer {
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

func poolForHubClients(timeoutInSec time.Duration, hub *Hub, username string, t *testing.T, isClientMore bool) {
	timeout := time.After(timeoutInSec * time.Second)
	tick := time.Tick(10 * time.Millisecond)

	for {
		select {
		case <-timeout:
			t.Fatalf("timeout waiting for client")
		case <-tick:
			if isClientMore {
				if len(hub.clients[username]) > 0 {
					return
				}
			} else {
				if len(hub.clients[username]) == 0 {
					return
				}
			}
		}
	}
}

func TestConnectionWithNoUserName(t *testing.T) {
	testStruct := testStruct(true, NewClient)
	defer testStruct.server.Close()

	url := "ws" + testStruct.server.URL[len("http"):]

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Errorf("websocket dial error: %v", err)
	}
	defer ws.Close()

	ws.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, _, err = ws.ReadMessage()
	if err == nil {
		t.Errorf("websocket connection should return an error")
	}

	closeErr, ok := err.(*websocket.CloseError)
	if !ok {
		t.Fatalf("expected CloseError, got: %T (%v)", err, err)
	}

	if closeErr.Text != UsernameDoesNotExist {
		t.Errorf("expected close reason %q, got %q", UsernameDoesNotExist, closeErr.Text)
	}
}

func TestHubClientOnConnect(t *testing.T) {
	testStruct := testStruct(true, NewClient)
	defer testStruct.server.Close()

	username := "erik"

	url := "ws" + testStruct.server.URL[len("http"):] + fmt.Sprintf("?username=%v", username)

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Errorf("websocket dial error: %v", err)
	}
	defer ws.Close()

	poolForHubClients(2, testStruct.hub, username, t, true)

	if len(testStruct.hub.clients[username]) != 1 {
		t.Errorf("hub client should be 1")
	}
}

func TestMessageSendAndReceive(t *testing.T) {
	testStruct := testStruct(true, NewClient)
	defer testStruct.server.Close()

	username := "erik"
	username2 := "hansen"

	url := "ws" + testStruct.server.URL[len("http"):] + fmt.Sprintf("?username=%v", username)
	url2 := "ws" + testStruct.server.URL[len("http"):] + fmt.Sprintf("?username=%v", username2)

	client, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Errorf("websocket dial error: %v", err)
	}
	defer client.Close()

	client2, _, err := websocket.DefaultDialer.Dial(url2, nil)
	if err != nil {
		t.Errorf("websocket dial error: %v", err)
	}
	defer client2.Close()

	poolForHubClients(2, testStruct.hub, username, t, true)
	poolForHubClients(2, testStruct.hub, username2, t, true)

	if len(testStruct.hub.clients[username]) != 1 {
		t.Errorf("hub client %v should be 1, but got: %v", username, len(testStruct.hub.clients[username]))
	}

	if len(testStruct.hub.clients[username2]) != 1 {
		t.Errorf("hub client %v should be 1, but got: %v", username2, len(testStruct.hub.clients[username2]))
	}

	client1Msg := []byte("hello")
	if err := client.WriteMessage(websocket.TextMessage, client1Msg); err != nil {
		t.Errorf("client write error: %v", err)
	}

	_, _, _ = client.ReadMessage() // Join message from user2
	_, client1ReadMsg, err := client.ReadMessage()
	if err != nil {
		t.Errorf("client1 read error: %v", err)
	}

	_, client2ReadMsg, err := client2.ReadMessage()
	if err != nil {
		t.Errorf("client2 read error: %v", err)
	}

	var jsonFromUser1Reader MessageJSON
	var jsonFromUser2Reader MessageJSON

	json.Unmarshal(client1ReadMsg, &jsonFromUser1Reader)
	json.Unmarshal(client2ReadMsg, &jsonFromUser2Reader)

	if jsonFromUser1Reader.Message != string(client1Msg) {
		t.Errorf("client2 expect: %v, but got %v", string(client1Msg), jsonFromUser1Reader.Message)
	}

	if jsonFromUser2Reader.Message != string(client1Msg) {
		t.Errorf("client2 expect: %v, but got %v", string(client1Msg), jsonFromUser2Reader.Message)
	}
}

func TestClientLeave(t *testing.T) {
	testStruct := testStruct(true, NewClient)
	defer testStruct.server.Close()
	username := "erik"

	url := "ws" + testStruct.server.URL[len("http"):] + fmt.Sprintf("?username=%v", username)
	client, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Errorf("websocket dial error: %v", err)
	}
	defer client.Close()
	poolForHubClients(2, testStruct.hub, username, t, true)

	if len(testStruct.hub.clients[username]) != 1 {
		t.Errorf("hub client %v should be 1, but got: %v", username, len(testStruct.hub.clients[username]))
	}

	client.NetConn().Close()
	poolForHubClients(2, testStruct.hub, username, t, false)

	if len(testStruct.hub.clients[username]) != 0 {
		t.Errorf("hub client %v should be 0, but got: %v", username, len(testStruct.hub.clients[username]))
	}
}

func TestIdleUsersShouldDisconnect(t *testing.T) {
	testStruct := testStruct(true, newTestClient)
	defer testStruct.server.Close()
	username := "erik"

	url := "ws" + testStruct.server.URL[len("http"):] + fmt.Sprintf("?username=%v", username)
	client, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Errorf("websocket dial error: %v", err)
	}
	defer client.Close()
	poolForHubClients(2, testStruct.hub, username, t, true)

	if len(testStruct.hub.clients[username]) != 1 {
		t.Errorf("Should have one client, but got %v", len(testStruct.hub.clients[username]))
	}

	poolForHubClients(10, testStruct.hub, username, t, false)

	if len(testStruct.hub.clients[username]) != 0 {
		t.Errorf("Should have 0 client, but got %v", len(testStruct.hub.clients[username]))
	}
}

func TestShouldBroadCastJoinAndQuitEvent(t *testing.T) {
	testStruct := testStruct(true, NewClient)
	defer testStruct.server.Close()
	username := "erik"
	username2 := "arne"

	url := "ws" + testStruct.server.URL[len("http"):] + fmt.Sprintf("?username=%v", username)
	url2 := "ws" + testStruct.server.URL[len("http"):] + fmt.Sprintf("?username=%v", username2)

	client, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Errorf("websocket dial error: %v", err)
	}
	defer client.Close()

	client2, _, err := websocket.DefaultDialer.Dial(url2, nil)
	if err != nil {
		t.Errorf("websocket dial error: %v", err)
	}
	defer client2.Close()

	_, msg, _ := client.ReadMessage()
	var jsonStructJoin MessageJSON
	json.Unmarshal(msg, &jsonStructJoin)

	if jsonStructJoin.Type != messageTypeJoin {
		t.Errorf("Should have type %v, but got %v", messageTypeJoin, jsonStructJoin.Type)
	}

	client2.Close()

	var jsonStructQuit MessageJSON
	_, msg, _ = client.ReadMessage()
	json.Unmarshal(msg, &jsonStructQuit)

	if jsonStructQuit.Type != messageTypeQuit {
		t.Errorf("Should have type %v, but got %v", messageTypeQuit, jsonStructQuit.Type)
	}

	if len(testStruct.hub.clients[username2]) != 0 {
		t.Errorf("Should have 0 client, but got %v", len(testStruct.hub.clients[username2]))
	}
}
