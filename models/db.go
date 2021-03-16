package models

import (
	"gorm.io/gorm"
)

var db *gorm.DB

// GetDB returns db
func GetDB() *gorm.DB {
	return db
}

// Init db using by models with conn
func Init(conn *gorm.DB) {
	db = conn
}

// Migrate ...
func Migrate() {
	GetDB().Debug().AutoMigrate(&Account{}, &Note{})
}
