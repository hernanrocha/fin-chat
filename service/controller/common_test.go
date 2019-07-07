package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hernanrocha/fin-chat/service/models"
)

func generateToken(t *testing.T, router *gin.Engine) string {
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

	// Perform Register
	w := performRequest(router, "POST", "/register", req)

	// Check
	require.Equal(t, http.StatusOK, w.Code)

	// Login Request
	req = gin.H{
		"username": fmt.Sprintf("user-%d", userID),
		"password": password,
	}

	// Perform Login
	w = performRequest(router, "POST", "/login", req)

	// Check
	require.Equal(t, http.StatusOK, w.Code)

	var resp gin.H
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	require.Nil(t, err)

	token, ok := resp["token"]
	require.True(t, ok)
	tokenStr, ok := token.(string)
	require.True(t, ok)

	return tokenStr
}

func assertUnauthorized(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp gin.H
	err := json.Unmarshal([]byte(w.Body.String()), &resp)

	require.Nil(t, err)
	code, ok := resp["code"]
	assert.True(t, ok)
	assert.EqualValues(t, http.StatusUnauthorized, code)
}

func performRequest(r http.Handler, method, path string, body gin.H) *httptest.ResponseRecorder {
	bodyRaw, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, path, bytes.NewReader(bodyRaw))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performAuthRequest(r http.Handler, method, path string, body gin.H, token string) *httptest.ResponseRecorder {
	bodyRaw, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, path, bytes.NewReader(bodyRaw))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func SetupDatabase() error {
	dbconn := getEnv("DB_CONNECTION", "host=localhost port=15432 user=postgres password=postgres dbname=finchat sslmode=disable")
	db, err := gorm.Open("postgres", dbconn)

	if err != nil {
		return err
	}

	models.Setup(db)

	return err
}
