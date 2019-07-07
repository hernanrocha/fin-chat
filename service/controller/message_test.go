package controller

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hernanrocha/fin-chat/service/viewmodels"
	"github.com/hernanrocha/fin-chat/service/hub/mocks"
)

func TestMessageListCreateGet(t *testing.T) {
	require.Nil(t, SetupDatabase())

	mockHub := mocks.NewMockHub()
	mockHub.On("BroadcastMessage", mock.AnythingOfType("viewmodels.MessageView")).
		Return().Once()
	router := SetupRouter(mockHub)

	token := generateToken(t, router)

	// Create Room
	roomName := fmt.Sprintf("Room %d", rand.Int())
	req := gin.H{
		"name": roomName,
	}
	w := performAuthRequest(router, "POST", "/api/v1/rooms", req, token)
	require.Equal(t, http.StatusOK, w.Code)

	var createRoomResp viewmodels.CreateRoomResponse
	err := json.Unmarshal([]byte(w.Body.String()), &createRoomResp)
	require.Nil(t, err)

	path := fmt.Sprintf("/api/v1/rooms/%d/messages", createRoomResp.ID)

	// List Messages
	w = performAuthRequest(router, "GET", path, nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp viewmodels.ListMessageResponse
	err = json.Unmarshal([]byte(w.Body.String()), &resp)
	require.Nil(t, err)
	assert.Empty(t, resp.Messages)

	// Create Message
	message := fmt.Sprintf("Message %d", rand.Int())
	req = gin.H{
		"text": message,
	}
	w = performAuthRequest(router, "POST", path, req, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var createResp viewmodels.CreateMessageResponse
	err = json.Unmarshal([]byte(w.Body.String()), &createResp)
	require.Nil(t, err)
	assert.NotZero(t, createResp.ID)
	assert.Equal(t, message, createResp.Text)

	// List Rooms (length should be one more)
	w = performAuthRequest(router, "GET", path, nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal([]byte(w.Body.String()), &resp)
	require.Nil(t, err)
	require.Equal(t, 1, len(resp.Messages))
}

func TestMessageUnauthorized(t *testing.T) {
	require.Nil(t, SetupDatabase())
	router := SetupRouter(nil)

	w := performRequest(router, "POST", "/api/v1/rooms/1/messages", nil)
	assertUnauthorized(t, w)

	w = performRequest(router, "GET", "/api/v1/rooms/1/messages", nil)
	assertUnauthorized(t, w)
}
