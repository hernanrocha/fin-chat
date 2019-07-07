package hub

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hernanrocha/fin-chat/service/viewmodels"
)

func TestAddRemoveClients(t *testing.T) {
	mock1 := NewMockMessageHandler()
	mock2 := NewMockMessageHandler()

	mock1.On("GetID").Return("mock1")
	mock2.On("GetID").Return("mock2")

	hub := NewHub()
	assert.Len(t, hub.clients, 0)

	hub.addClient(mock1)
	assert.Len(t, hub.clients, 1)

	hub.addClient(mock2)
	assert.Len(t, hub.clients, 2)

	hub.removeClient(mock1)
	assert.Len(t, hub.clients, 1)

	hub.removeClient(mock2)
	assert.Len(t, hub.clients, 0)

	mock1.AssertExpectations(t)
	mock2.AssertExpectations(t)
}

func TestSendMessage(t *testing.T) {
	mock1 := NewMockMessageHandler()
	mock2 := NewMockMessageHandler()

	msg := viewmodels.MessageView{
		ID:        2,
		Text:      "Text",
		Username:  "Username",
		CreatedAt: time.Now(),
		RoomID:    3,
	}

	f := func(args mock.Arguments) {
		arg, ok := args.Get(0).(viewmodels.MessageView)
		require.True(t, ok)
		require.Equal(t, msg, arg)
	}

	mock1.On("GetID").Return("mock1")
	mock2.On("GetID").Return("mock2")
	mock1.On("HandleMessage", mock.AnythingOfType("viewmodels.MessageView")).
		Run(f).Return(nil).Once()
	mock2.On("HandleMessage", mock.AnythingOfType("viewmodels.MessageView")).
		Run(f).Return(nil).Once()

	hub := NewHub()
	hub.addClient(mock1)
	hub.addClient(mock2)
	hub.broadcastMessage(msg)

	mock1.AssertExpectations(t)
	mock2.AssertExpectations(t)
}

type MockMessageHandler struct {
	mock.Mock
}

func NewMockMessageHandler() *MockMessageHandler {
	return &MockMessageHandler{}
}

func (h *MockMessageHandler) GetID() string {
	args := h.Called()
	return args.String(0)
}

func (h *MockMessageHandler) HandleMessage(msg viewmodels.MessageView) error {
	args := h.Called(msg)
	return args.Error(0)
}