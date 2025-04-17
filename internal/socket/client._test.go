package socket

import (
	"encoding/json"
)

func fromByte(b []byte) (*MessageJSON, error) {
	var msg MessageJSON
	err := json.Unmarshal(b, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

//func TestMessageSendAndReceiveWhenOnJoinAndLeave(t *testing.T) {
//	serverHub := testServerHub(true, newClient)
//	defer serverHub.server.Close()
//
//	username := "erik"
//	username2 := "hansen"
//	channel := "test"
//
//	client, _, err := websocket.DefaultDialer.Dial(serverHub.wsUrl(username, channel), nil)
//	if err != nil {
//		t.Errorf("websocket dial error: %v", err)
//	}
//	defer client.Close()
//
//	client2, _, err := websocket.DefaultDialer.Dial(serverHub.wsUrl(username2, channel), nil)
//	if err != nil {
//		t.Errorf("websocket dial error: %v", err)
//	}
//	defer client2.Close()
//
//	err = pollingForClientInChannels(pollingForClientInChannelsArgs{
//		sec:     2,
//		hub:     serverHub.hub,
//		channel: channel,
//		cap:     2,
//	})
//	if err != nil {
//		t.Errorf("pollingForClientInChannels error: %v", err)
//	}
//
//	channelLength := len(serverHub.hub.channels[channel])
//	if channelLength != 2 {
//		t.Errorf("should be 2 client in channel %s but got %d", channel, channelLength)
//	}
//
//	_ = client.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
//	_ = client2.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
//
//	_, msg, _ := client.ReadMessage()
//	joinMsg, err := fromByte(msg)
//	if err != nil {
//		t.Errorf("json unmarshal error: %v", err)
//	}
//
//	if joinMsg.Event != messageTypeJoin && joinMsg.Username != username2 {
//		t.Errorf("should be join event, but got %s from user %s but got %s", joinMsg.Event, username2, joinMsg.Username)
//	}
//
//	client1Msg := []byte("hello")
//	if err := client.WriteMessage(websocket.TextMessage, client1Msg); err != nil {
//		t.Errorf("client write error: %v", err)
//	}
//
//	_, msg, _ = client.ReadMessage()
//	helloMsg, err := fromByte(msg)
//	if err != nil {
//		t.Errorf("json unmarshal error: %v", err)
//	}
//
//	if helloMsg.Message != string(client1Msg) {
//		t.Errorf("should receive own message %s but got %s", client1Msg, helloMsg.Message)
//	}
//
//	_, msg, _ = client2.ReadMessage()
//	helloMsg, err = fromByte(msg)
//	if err != nil {
//		t.Errorf("json unmarshal error: %v", err)
//	}
//
//	if helloMsg.Message != string(client1Msg) {
//		t.Errorf("%s should receive message %s, but got %s", username2, string(client1Msg), helloMsg.Message)
//	}
//
//	_ = client2.Close()
//
//	err = pollingForClientInChannels(pollingForClientInChannelsArgs{
//		sec:     2,
//		hub:     serverHub.hub,
//		channel: channel,
//		cap:     1,
//	})
//	if err != nil {
//		t.Errorf("pollingForClientInChannels error: %v", err)
//	}
//
//	_, msg, _ = client.ReadMessage()
//	quitMsg, err := fromByte(msg)
//	if err != nil {
//		t.Errorf("json unmarshal error: %v", err)
//	}
//
//	if quitMsg.Event != messageTypeQuit && quitMsg.Username != username2 {
//		t.Errorf("should leave event type %s, but got %s. Username should be %s, but got %s", messageTypeQuit, quitMsg.Event, username2, quitMsg.Username)
//	}
//}
//
//func TestIdleUsersShouldDisconnect(t *testing.T) {
//	serverHub := testServerHub(true, newTestClient)
//	defer serverHub.server.Close()
//	username := "erik"
//	channel := "test"
//
//	client, _, err := websocket.DefaultDialer.Dial(serverHub.wsUrl(username, channel), nil)
//	if err != nil {
//		t.Errorf("websocket dial error: %v", err)
//	}
//	defer client.Close()
//
//	err = pollingForClientInChannels(pollingForClientInChannelsArgs{
//		sec:     2,
//		hub:     serverHub.hub,
//		channel: channel,
//		cap:     1,
//	})
//	if err != nil {
//		t.Errorf("pollingForClientInChannels error: %v", err)
//	}
//
//	if len(serverHub.hub.channels[channel]) != 1 {
//		t.Errorf("Should have one client, but got %v", len(serverHub.hub.channels[channel]))
//	}
//
//	err = pollingForClientInChannels(pollingForClientInChannelsArgs{
//		sec:     2,
//		hub:     serverHub.hub,
//		channel: channel,
//		cap:     0,
//	})
//
//	if err != nil {
//		t.Errorf("pollingForClientInChannels error: %v", err)
//	}
//
//	if len(serverHub.hub.channels[channel]) != 0 {
//		t.Errorf("Should have 0 client, but got %v", len(serverHub.hub.channels[channel]))
//	}
//}
