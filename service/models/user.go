package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"type:varchar(100);unique_index"`
	Password  string
	Email     string `gorm:"type:varchar(100);unique_index"`
	FirstName string
	LastName  string
}
