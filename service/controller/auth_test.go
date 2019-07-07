package controller

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterLogin(t *testing.T) {
	require.Nil(t, SetupDatabase())
	router := SetupRouter(nil)

	rand.Seed(int64(time.Now().Nanosecond()))
	userID := rand.Int()
	password := "password"

	// Register Request
	req := gin.H{
		"username":   fmt.Sprintf("user-%d", userID),
		"email":      fmt.Sprintf("user-%d@mail.com", userID),
		"password":   password,
		"first_name": "First Name",
		"last_name":  "Last Name",
	}

	// Register Response
	resp := gin.H{
		"username":   req["username"],
		"email":      req["email"],
		"first_name": req["first_name"],
		"last_name":  req["last_name"],
	}

	// Perform Register
	w := performRequest(router, "POST", "/register", req)

	// Check
	assert.Equal(t, http.StatusOK, w.Code)

	var actualResp gin.H
	err := json.Unmarshal([]byte(w.Body.String()), &actualResp)
	require.Nil(t, err)
	assert.Equal(t, resp, actualResp)

	// Login Request
	loginReq := gin.H{
		"username": req["username"],
		"password": password,
	}

	// Perform Login
	w = performRequest(router, "POST", "/login", loginReq)

	// Check
	assert.Equal(t, http.StatusOK, w.Code)

	var actualLoginResp gin.H
	err = json.Unmarshal([]byte(w.Body.String()), &actualLoginResp)
	require.Nil(t, err)

	code, ok := actualLoginResp["code"]
	assert.True(t, ok)
	assert.EqualValues(t, http.StatusOK, code)

	expire, ok := actualLoginResp["expire"]
	assert.True(t, ok)
	expireStr, ok := expire.(string)
	assert.True(t, ok)
	expireTime, err := time.Parse(time.RFC3339, expireStr)
	assert.Nil(t, err)
	assert.True(t, expireTime.After(time.Now()), fmt.Sprintf("Expire time should be in the future. Actual: %s", expireTime))

	token, ok := actualLoginResp["token"]
	assert.True(t, ok)
	_, ok = token.(string)
	assert.True(t, ok)
}

func TestRegisterErrorNoPassword(t *testing.T) {
	require.Nil(t, SetupDatabase())
	router := SetupRouter(nil)

	rand.Seed(int64(time.Now().Nanosecond()))
	userID := rand.Int()

	// Register Request (without required field password)
	req := gin.H{
		"username": fmt.Sprintf("user-%d", userID),
		"email":    fmt.Sprintf("user-%d@mail.com", userID),
	}

	// Perform Register
	w := performRequest(router, "POST", "/register", req)

	// Check
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp gin.H
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	require.Nil(t, err)

	e, ok := resp["error"]
	assert.True(t, ok)
	eStr, ok := e.(string)
	assert.True(t, ok)
	assert.Contains(t, eStr, "Password")
}

func TestLoginInvalidCredentials(t *testing.T) {
	require.Nil(t, SetupDatabase())
	router := SetupRouter(nil)

	// Login Request (with invalid credentials)
	req := gin.H{
		"username": "invalid-user",
		"password": "invalid-password",
	}

	// Perform Login
	w := performRequest(router, "POST", "/login", req)

	// Check
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp gin.H
	err := json.Unmarshal([]byte(w.Body.String()), &resp)

	require.Nil(t, err)
	code, ok := resp["code"]
	assert.True(t, ok)
	assert.EqualValues(t, http.StatusUnauthorized, code)

	message, ok := resp["message"]
	assert.True(t, ok)
	messageStr, ok := message.(string)
	assert.True(t, ok)
	assert.Equal(t, "incorrect Username or Password", messageStr)
}
