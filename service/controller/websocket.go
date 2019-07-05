package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func wshandler(w http.ResponseWriter, r *http.Request, h *Hub) {
	fmt.Println("NEW WEBSOCKET")
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}

	h.AddClientChan <- conn

	for {
		_, _, err = conn.ReadMessage()
		if err != nil {
			h.RemoveClientChan <- conn
			return
		}
	}
}

// WebSocketController ...
type WebSocketController struct {
	hub *Hub
}

// NewWebSocketController ...
func NewWebSocketController(hub *Hub) *WebSocketController {
	return &WebSocketController{
		hub: hub,
	}
}

func (c *WebSocketController) WebSocket(ctx *gin.Context) {
	wshandler(ctx.Writer, ctx.Request, c.hub)
}
