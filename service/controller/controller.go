package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hernanrocha/fin-chat/service/models"
)

// Controller example
type Controller struct {
}

// NewController example
func NewController() *Controller {
	return &Controller{}
}

type PingResponse struct {
	Message string `json:"message"`
}

type UserView struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RoomView struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type RegisterRequest struct {
	UserView
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	UserView
}

type ListRoomResponse struct {
	Rooms []RoomView `json:"rooms"`
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}

type CreateRoomResponse struct {
	RoomView
}

type GetRoomResponse struct {
	RoomView
}

type ListMessageResponse struct {
	Messages []MessageView `json:"messages"`
}

// Ping godoc
// @Summary Simple Ping/Pong protocol
// @Description Receive a Ping request and send a Pong response
// @Produce  json
// @Success 200 {object} controller.PingResponse
// @Router /ping [get]
func (c *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, PingResponse{Message: "pong"})
}

// Register User
// @Summary Register User
// @Description Register User in database
// @Produce  json
// @Success 200 {object} controller.RegisterResponse
// @Router /register [post]
func (c *Controller) Register(ctx *gin.Context) {
	var json RegisterRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := models.GetDB()

	user := &models.User{
		Username:  json.Username,
		Password:  json.Password,
		Email:     json.Email,
		FirstName: json.FirstName,
		LastName:  json.LastName,
	}

	if err := db.Create(user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := &RegisterResponse{
		UserView{
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

// List Rooms
// @Summary List Rooms
// @Description List Rooms in database
// @Produce  json
// @Success 200 {object} controller.ListRoomResponse
// @Router /rooms [get]
func (c *Controller) ListRooms(ctx *gin.Context) {
	db := models.GetDB()

	var rooms []models.Room
	if err := db.Find(&rooms).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomList := make([]RoomView, len(rooms))
	for i, r := range rooms {
		roomList[i] = RoomView{
			ID:   r.ID,
			Name: r.Name,
		}
	}

	response := &ListRoomResponse{
		Rooms: roomList,
	}

	ctx.JSON(http.StatusOK, response)
}

// CreateRoom
// @Summary Create Room
// @Description Create Room in database
// @Produce  json
// @Success 200 {object} controller.CreateRoomResponse
// @Router /rooms [post]
func (c *Controller) CreateRoom(ctx *gin.Context) {
	var json CreateRoomRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := models.GetDB()

	room := &models.Room{
		Name: json.Name,
	}

	if err := db.Create(room).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := &CreateRoomResponse{
		RoomView{
			ID:   room.ID,
			Name: room.Name,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

// GetRoom
// @Summary Get Room
// @Description Get Room by ID
// @Produce  json
// @Success 200 {object} controller.GetRoomResponse
// @Router /rooms/:id [get]
func (c *Controller) GetRoom(ctx *gin.Context) {
	db := models.GetDB()

	id := ctx.Params.ByName("id")
	var room models.Room

	if err := db.Where("id = ?", id).Find(&room).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := &GetRoomResponse{
		RoomView{
			ID:   room.ID,
			Name: room.Name,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

type MessageView struct {
	ID     uint   `json:"id"`
	Text   string `json:"text"`
	UserID uint   `json:"user_id"`
	RoomID uint   `json:"room_id"`
}

type CreateMessageRequest struct {
	Text   string `json:"text"`
	UserID uint   `json:"user_id"`
	RoomID uint   `json:"room_id"`
}

type CreateMessageResponse struct {
	MessageView
}

// CreateMessage
// @Summary Create Message
// @Description Create Message in database
// @Produce  json
// @Success 200 {object} controller.CreateMessageResponse
// @Router /messages [post]
func (c *Controller) CreateMessage(ctx *gin.Context) {
	var json CreateMessageRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := models.GetDB()

	message := &models.Message{
		Text:   json.Text,
		RoomID: json.RoomID,
		UserID: json.UserID,
	}

	if err := db.Create(message).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := &CreateMessageResponse{
		MessageView{
			ID:     message.ID,
			Text:   message.Text,
			RoomID: message.RoomID,
			UserID: message.UserID,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

// List Room Messages
// @Summary List Room Messages
// @Description List last Room Messages in database
// @Produce  json
// @Success 200 {object} controller.ListMessageResponse
// @Router /rooms/:id/messages [get]
func (c *Controller) ListRoomMessages(ctx *gin.Context) {
	db := models.GetDB()

	var messages []models.Message
	id := ctx.Params.ByName("id")

	if err := db.Where("room_id = ?", id).Limit(50).Order("created_at asc").Find(&messages).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageList := make([]MessageView, len(messages))
	for i, m := range messages {
		messageList[i] = MessageView{
			ID:     m.ID,
			Text:   m.Text,
			UserID: m.UserID,
			RoomID: m.RoomID,
		}
	}

	response := &ListMessageResponse{
		Messages: messageList,
	}

	ctx.JSON(http.StatusOK, response)
}
