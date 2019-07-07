package hub

import (
	"log"

	"github.com/hernanrocha/fin-chat/service/viewmodels"
)

type MessageHandler interface {
	GetID() string
	HandleMessage(msg viewmodels.MessageView) error
}

type HubInterface interface {
	AddClient(h MessageHandler)
	RemoveClient(h MessageHandler)
	BroadcastMessage(m viewmodels.MessageView)
}

type Hub struct {
	clients          map[string]MessageHandler
	AddClientChan    chan MessageHandler
	RemoveClientChan chan MessageHandler
	BroadcastChan    chan viewmodels.MessageView
}

func NewHub() *Hub {
	return &Hub{
		clients:          make(map[string]MessageHandler),
		AddClientChan:    make(chan MessageHandler),
		RemoveClientChan: make(chan MessageHandler),
		BroadcastChan:    make(chan viewmodels.MessageView),
	}
}

func (h *Hub) Run() {
	go h.run()
}

func (h *Hub) RemoveClient(handler MessageHandler) {
	h.RemoveClientChan <- handler
}

func (h *Hub) AddClient(handler MessageHandler) {
	h.AddClientChan <- handler
}

func (h *Hub) BroadcastMessage(m viewmodels.MessageView) {
	h.BroadcastChan <- m
}

func (h *Hub) run() {
	for {
		select {
		case handler := <-h.AddClientChan:
			h.addClient(handler)
		case handler := <-h.RemoveClientChan:
			h.removeClient(handler)
		case m := <-h.BroadcastChan:
			h.broadcastMessage(m)
		}
	}
}

func (h *Hub) addClient(handler MessageHandler) {
	log.Println("Adding client...")
	h.clients[handler.GetID()] = handler
}

func (h *Hub) removeClient(handler MessageHandler) {
	log.Println("Removing client...")
	delete(h.clients, handler.GetID())
}

func (h *Hub) broadcastMessage(msg viewmodels.MessageView) {
	log.Printf("Broadcasting message: %s\n", msg.Text)
	for _, handler := range h.clients {
		if err := handler.HandleMessage(msg); err != nil {
			log.Printf("Error broadcasting message: %s\n", err)
			delete(h.clients, handler.GetID())
		}
	}
}
