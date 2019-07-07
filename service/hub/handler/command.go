package handler

import (
	"log"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/hernanrocha/fin-chat/messenger"
	"github.com/hernanrocha/fin-chat/service/models"
	"github.com/hernanrocha/fin-chat/service/viewmodels"
	"github.com/hernanrocha/fin-chat/service/hub"
)

type CmdMessageHandler struct {
	ID   string
	msg  messenger.BotCommandMessenger
	db   *gorm.DB
	hub  hub.HubInterface
	user models.User
}

func NewCmdMessageHandler(ID string, msg messenger.BotCommandMessenger, hub hub.HubInterface, db *gorm.DB) (*CmdMessageHandler, error) {
	handler := &CmdMessageHandler{
		ID:  ID,
		msg: msg,
		hub: hub,
		db:  db,
	}

	if err := handler.setup(); err != nil {
		return nil, err
	}

	return handler, nil
}

func (h *CmdMessageHandler) HandleMessage(msg viewmodels.MessageView) error {
	if strings.HasPrefix(msg.Text, "/stock=") {
		cmd := msg.Text[7:]
		log.Printf("Sending command '%s' to StockBot...\n", cmd)
		h.msg.Publish(strconv.Itoa(int(msg.RoomID)), cmd)
	}

	return nil
}

func (h *CmdMessageHandler) GetID() string {
	return h.ID
}

func (h *CmdMessageHandler) CmdResponseHandler(botMsg messenger.BotMessage) error {
	message := &models.Message{
		Text:   botMsg.Message,
		RoomID: botMsg.RoomID,
		UserID: h.user.ID,
	}

	if err := h.db.Create(message).Error; err != nil {
		log.Println("Error creating new message from bot: ", err)
		return err
	}

	mv := viewmodels.MessageView{
		ID:        message.ID,
		Text:      message.Text,
		RoomID:    message.RoomID,
		Username:  h.user.Username,
		CreatedAt: message.CreatedAt,
	}

	// Broadcast message to Hub
	h.hub.BroadcastMessage(mv)

	return nil
}

func (h *CmdMessageHandler) setup() error {
	h.user = models.User{
		Username: "Bot",
		Email:    "bot@mail.com",
	}

	if err := h.db.Where("username = ?", h.user.Username).First(&h.user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println("Creating bot user")
			if err := h.db.Create(&h.user).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
