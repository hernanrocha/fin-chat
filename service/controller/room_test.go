package controller

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hernanrocha/fin-chat/service/viewmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoomListCreateGet(t *testing.T) {
	require.Nil(t, SetupDatabase())
	router := SetupRouter(nil)

	token := generateToken(t, router)

	// List Rooms
	w := performAuthRequest(router, "GET", "/api/v1/rooms", nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp viewmodels.ListRoomResponse
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	require.Nil(t, err)
	oldCount := len(resp.Rooms)

	// Create Room
	roomName := fmt.Sprintf("Room %d", rand.Int())
	req := gin.H{
		"name": roomName,
	}
	w = performAuthRequest(router, "POST", "/api/v1/rooms", req, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var createResp viewmodels.CreateRoomResponse
	err = json.Unmarshal([]byte(w.Body.String()), &createResp)
	require.Nil(t, err)
	assert.Equal(t, roomName, createResp.Name)

	// Get Room by ID
	path := fmt.Sprintf("/api/v1/rooms/%d", createResp.ID)
	w = performAuthRequest(router, "GET", path, req, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var getResp viewmodels.GetRoomResponse
	err = json.Unmarshal([]byte(w.Body.String()), &getResp)
	require.Nil(t, err)
	assert.EqualValues(t, createResp, getResp)

	// List Rooms (length should be one more)
	w = performAuthRequest(router, "GET", "/api/v1/rooms", nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal([]byte(w.Body.String()), &resp)
	require.Nil(t, err)
	assert.Equal(t, oldCount+1, len(resp.Rooms))
}

func TestRoomsUnauthorized(t *testing.T) {
	require.Nil(t, SetupDatabase())
	router := SetupRouter(nil)

	w := performRequest(router, "GET", "/api/v1/rooms", nil)
	assertUnauthorized(t, w)

	w = performRequest(router, "POST", "/api/v1/rooms", nil)
	assertUnauthorized(t, w)

	w = performRequest(router, "GET", "/api/v1/rooms/1", nil)
	assertUnauthorized(t, w)
}

func TestGetRoomNotExist(t *testing.T) {
	require.Nil(t, SetupDatabase())
	router := SetupRouter(nil)

	token := generateToken(t, router)

	w := performAuthRequest(router, "GET", "/api/v1/rooms/1111111", nil, token)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateRoomInvalidRequest(t *testing.T) {
	require.Nil(t, SetupDatabase())
	router := SetupRouter(nil)

	token := generateToken(t, router)

	// Create Room (without required field name)
	req := gin.H{}
	w := performAuthRequest(router, "POST", "/api/v1/rooms", req, token)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
