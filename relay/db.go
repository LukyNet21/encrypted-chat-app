package main

import (
	"encrypted-chat-relay/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{}, &models.Message{})

	return db
}
