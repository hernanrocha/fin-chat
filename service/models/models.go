package models

import (
	"github.com/jinzhu/gorm"
)

// DB Instance
var DB *gorm.DB

// Setup database and run migrations
func Setup(db *gorm.DB) error {
	// Migrate User
	if err := db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}

	// Migrate Room
	if err := db.AutoMigrate(&Room{}).Error; err != nil {
		return err
	}

	// Migrate Message
	db = db.AutoMigrate(&Message{}).
		AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").
		AddForeignKey("room_id", "rooms(id)", "CASCADE", "CASCADE")

	if db.Error != nil {
		return db.Error
	}

	// Store DB
	DB = db
	return nil
}

// GetDB Get database instance
func GetDB() *gorm.DB {
	return DB
}
