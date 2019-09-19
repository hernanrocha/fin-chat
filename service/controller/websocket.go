package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/hernanrocha/fin-chat/service/hub"
	"github.com/hernanrocha/fin-chat/service/hub/handler"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// WebSocketController ...
type WebSocketController struct {
	hub hub.HubInterface
}

// NewWebSocketController ...
func NewWebSocketController(hub hub.HubInterface) *WebSocketController {
	return &WebSocketController{
		hub: hub,
	}
}

func (c *WebSocketController) WebSocket(ctx *gin.Context) {
	log.Println("Creating new WebSocket")
	conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %s\n", err)
		return
	}

	handler := handler.NewWebSocketMessageHandler(conn)
	c.hub.AddClient(handler)

	for {
		_, _, err = conn.ReadMessage()
		if err != nil {
			log.Println("ERROR ON WEBSOCKET: CLOSING...")
			log.Println(err)
			c.hub.RemoveClient(handler)
			return
		}
	}
}
