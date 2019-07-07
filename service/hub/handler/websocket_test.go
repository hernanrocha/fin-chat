package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/hernanrocha/fin-chat/service/viewmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func TestWebSocketMessageHandler(t *testing.T) {
	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	wsh := NewWebSocketMessageHandler(ws)
	require.NotNil(t, wsh)
	assert.NotEmpty(t, wsh.GetID())

	msg := viewmodels.MessageView{
		RoomID:   100,
		Text:     "Sample message",
		Username: "BotUsername",
	}
	err = wsh.HandleMessage(msg)
	require.Nil(t, err)

	_, wsMsg, err := ws.ReadMessage()
	require.Nil(t, err)

	var msgReply viewmodels.MessageView
	err = json.Unmarshal(wsMsg, &msgReply)
	assert.Equal(t, msg, msgReply)
}
