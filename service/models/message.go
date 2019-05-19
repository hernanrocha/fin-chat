package models

import (
	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	Text   string
	UserID uint
	RoomID uint

	User *User
}
