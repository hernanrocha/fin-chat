package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hernanrocha/fin-chat/service/models"
)

func performRequest(r http.Handler, method, path string, body gin.H) *httptest.ResponseRecorder {
	bodyRaw, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, path, bytes.NewReader(bodyRaw))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func SetupDatabase() error {
	dbconn := "host=localhost port=15432 user=postgres password=postgres dbname=finchat sslmode=disable"
	db, err := gorm.Open("postgres", dbconn)

	if err != nil {
		return err
	}

	models.Setup(db)

	return err
}

func TestRegister(t *testing.T) {
	require.Nil(t, SetupDatabase())

	// Build our expected body
	body := gin.H{
		"username": "UserName",
	}

	router := SetupRouter(nil, nil)

	// Perform a GET request with that handler.
	w := performRequest(router, "POST", "/register", gin.H{
		"username":   "UserName",
		"email":      "UserName@gmail.com",
		"password":   "password",
		"first_name": "Hernan",
	})

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Convert the JSON response to a map
	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	require.Nil(t, err)

	value, exists := response["username"]
	assert.True(t, exists)
	assert.Equal(t, body["username"], value)
}
