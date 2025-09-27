package main

import (
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func InitializeDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("storage.db"))

	if err != nil {
		log.Fatal("Cannot connect to the database: ", err)
	}

	db.AutoMigrate(&User{}, &Book{})

	return db
}
