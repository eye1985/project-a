package socket

//func TestConnectionWithNoUserName(t *testing.T) {
//	testStruct := testServerHub(true, newClient)
//	defer testStruct.server.Close()
//
//	url := "ws" + testStruct.server.URL[len("http"):]
//
//	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
//	if err != nil {
//		t.Errorf("websocket dial error: %v", err)
//	}
//	defer ws.Close()
//
//	ws.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
//	_, _, err = ws.ReadMessage()
//	if err == nil {
//		t.Errorf("websocket connection should return an error")
//	}
//
//	closeErr, ok := err.(*websocket.CloseError)
//	if !ok {
//		t.Fatalf("expected CloseError, got: %T (%v)", err, err)
//	}
//
//	if closeErr.Text != usernameDoesNotExist {
//		t.Errorf("expected close reason %q, got %q", usernameDoesNotExist, closeErr.Text)
//	}
//}
//
//func TestConnectionWithNoChannelName(t *testing.T) {
//	testStruct := testServerHub(true, newClient)
//	defer testStruct.server.Close()
//
//	url := "ws" + testStruct.server.URL[len("http"):] + "?username=erik"
//
//	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
//	if err != nil {
//		t.Errorf("websocket dial error: %v", err)
//	}
//	defer ws.Close()
//
//	ws.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
//	_, _, err = ws.ReadMessage()
//	if err == nil {
//		t.Errorf("websocket connection should return an error")
//	}
//
//	closeErr, ok := err.(*websocket.CloseError)
//	if !ok {
//		t.Fatalf("expected CloseError, got: %T (%v)", err, err)
//	}
//
//	if closeErr.Text != channelNameDoesNotExist {
//		t.Errorf("expected close reason %q, got %q", channelNameDoesNotExist, closeErr.Text)
//	}
//}
