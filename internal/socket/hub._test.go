package socket

//func TestUpdateCorrectClientOnJoinLeave(t *testing.T) {
//	serverHub := testServerHub(false, newClient)
//	defer serverHub.server.Close()
//	channel := "testChannel"
//	username := "erik"
//	username2 := "per"
//
//	_, _, err := websocket.DefaultDialer.Dial(serverHub.wsUrl(username, channel), nil)
//	if err != nil {
//		t.Errorf("websocket dial error: %v", err)
//	}
//	client2, _, err := websocket.DefaultDialer.Dial(serverHub.wsUrl(username2, channel), nil)
//	if err != nil {
//		t.Errorf("websocket dial error: %v", err)
//	}
//
//	err = pollingForClientInChannels(pollingForClientInChannelsArgs{
//		sec:     2,
//		hub:     serverHub.hub,
//		channel: channel,
//		cap:     2,
//	})
//	if err != nil {
//		t.Errorf("pooling failed: %v", err)
//	}
//
//	res := []string{}
//	for _, c := range serverHub.hub.channels[channel] {
//		if c.username == username || c.username == username2 {
//			res = append(res, c.username)
//		}
//	}
//
//	if len(res) != 2 {
//		t.Errorf("clients on websocket should have two clients but got %d", len(res))
//	}
//
//	_ = client2.Close()
//	err = pollingForClientInChannels(pollingForClientInChannelsArgs{
//		sec:     2,
//		hub:     serverHub.hub,
//		channel: channel,
//		cap:     1,
//	})
//	if err != nil {
//		t.Errorf("pooling failed: %v", err)
//	}
//
//	if len(serverHub.hub.channels[channel]) != 1 {
//		t.Errorf("clients on websocket should have one clients but got %d", len(serverHub.hub.channels[channel]))
//	}
//}
