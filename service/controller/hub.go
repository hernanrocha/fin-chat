package controller

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/hernanrocha/fin-chat/service/viewmodels"
)

type Hub struct {
	clients          map[string]*websocket.Conn
	AddClientChan    chan *websocket.Conn
	RemoveClientChan chan *websocket.Conn
	BroadcastChan    chan viewmodels.MessageView
}

func NewHub() *Hub {
	return &Hub{
		clients:          make(map[string]*websocket.Conn),
		AddClientChan:    make(chan *websocket.Conn),
		RemoveClientChan: make(chan *websocket.Conn),
		BroadcastChan:    make(chan viewmodels.MessageView),
	}
}

func (h *Hub) Run() {
	go func() {
		for {
			select {
			case conn := <-h.AddClientChan:
				h.addClient(conn)
			case conn := <-h.RemoveClientChan:
				h.removeClient(conn)
			case m := <-h.BroadcastChan:
				h.broadcastMessage(m)
			}
		}
	}()
}

func (h *Hub) removeClient(conn *websocket.Conn) {
	fmt.Println("Removing client...")
	delete(h.clients, conn.RemoteAddr().String())
}

func (h *Hub) addClient(conn *websocket.Conn) {
	fmt.Println("Adding client...")
	h.clients[conn.RemoteAddr().String()] = conn
}

func (h *Hub) broadcastMessage(m viewmodels.MessageView) {
	fmt.Println("Broadcasting message: ", m.Text)
	for _, conn := range h.clients {
		err := websocket.WriteJSON(conn, m)
		if err != nil {
			fmt.Println("Error broadcasting message: ", err)
			h.removeClient(conn)
		}
	}
}
