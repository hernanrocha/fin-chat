package handler

import (
	"github.com/gorilla/websocket"
	"github.com/hernanrocha/fin-chat/service/viewmodels"
)

type WebSocketMessageHandler struct {
	ws *websocket.Conn
}

func NewWebSocketMessageHandler(ws *websocket.Conn) *WebSocketMessageHandler {
	return &WebSocketMessageHandler{
		ws: ws,
	}
}

func (h *WebSocketMessageHandler) HandleMessage(msg viewmodels.MessageView) error {
	return websocket.WriteJSON(h.ws, msg)
}

func (h *WebSocketMessageHandler) GetID() string {
	return h.ws.RemoteAddr().String()
}
