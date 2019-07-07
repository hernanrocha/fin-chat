package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/hernanrocha/fin-chat/service/hub"
	"github.com/hernanrocha/fin-chat/service/viewmodels"
)

type MockHub struct {
	mock.Mock
}

func NewMockHub() *MockHub {
	return &MockHub{}
}

func (hub *MockHub) AddClient(h hub.MessageHandler) {
	hub.Called(h)
}

func (hub *MockHub) RemoveClient(h hub.MessageHandler) {
	hub.Called(h)
}

func (hub *MockHub) BroadcastMessage(m viewmodels.MessageView) {
	hub.Called(m)
}
