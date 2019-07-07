package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hernanrocha/fin-chat/service/hub"
	"github.com/hernanrocha/fin-chat/service/models"
	"github.com/hernanrocha/fin-chat/service/viewmodels"
	"github.com/jinzhu/gorm"
)

// MessageController ...
type MessageController struct {
	hub hub.HubInterface
	db  *gorm.DB
}

// NewMessageController ...
func NewMessageController(hub hub.HubInterface) *MessageController {
	return &MessageController{
		hub: hub,
		db:  models.GetDB(),
	}
}

// CreateMessage godoc
// @Summary Create Message
// @Description Create Message in database
// @Tags Messages
// @Param Authorization header string true "JWT Token"
// @Param id path int true "Room ID"
// @Param user body viewmodels.CreateMessageRequest true "Message Data"
// @Produce  json
// @Success 200 {object} viewmodels.CreateMessageResponse
// @Router /api/v1/rooms/{id}/messages [post]
func (c *MessageController) CreateMessage(ctx *gin.Context) {
	var json viewmodels.CreateMessageRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userView, _ := ctx.Get("username")
	id := ctx.Params.ByName("id")
	uid, _ := strconv.Atoi(id)
	var user models.User
	if err := c.db.Where("username = ?", userView.(*viewmodels.UserView).Username).Find(&user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := &models.Message{
		Text:   json.Text,
		RoomID: uint(uid),
		UserID: user.ID,
	}

	if err := c.db.Create(message).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mv := viewmodels.MessageView{
		ID:        message.ID,
		Text:      message.Text,
		RoomID:    message.RoomID,
		Username:  user.Username,
		CreatedAt: message.CreatedAt,
	}

	// Broadcast message to Hub
	c.hub.BroadcastMessage(mv)

	response := &viewmodels.CreateMessageResponse{
		MessageView: mv,
	}

	ctx.JSON(http.StatusOK, response)
}

// ListRoomMessages godoc
// @Summary List Room Messages
// @Description List last Room Messages in database
// @Tags Messages
// @Param Authorization header string true "JWT Token"
// @Param id path int true "Room ID"
// @Produce  json
// @Success 200 {object} viewmodels.ListMessageResponse
// @Router /api/v1/rooms/{id}/messages [get]
func (c *MessageController) ListRoomMessages(ctx *gin.Context) {
	var messages []models.Message
	id := ctx.Params.ByName("id")

	db := c.db.Where("room_id = ?", id).Preload("User").Limit(50).Order("created_at desc").Find(&messages)
	if err := db.Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageList := make([]viewmodels.MessageView, len(messages))
	for i, m := range messages {
		messageList[i] = viewmodels.MessageView{
			ID:        m.ID,
			Text:      m.Text,
			RoomID:    m.RoomID,
			Username:  m.User.Username,
			CreatedAt: m.CreatedAt,
		}
	}

	response := &viewmodels.ListMessageResponse{
		Messages: messageList,
	}

	ctx.JSON(http.StatusOK, response)
}
