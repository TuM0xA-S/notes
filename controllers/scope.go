package controllers

import (
	"net/http"
	"notes/models"

	"gorm.io/gorm"
)

// useless because of small logic

// Published ...
func Published(db *gorm.DB) *gorm.DB {
	return db.Where("published")
}

// OwnedBy user from request
func OwnedBy(req *http.Request) func(db *gorm.DB) *gorm.DB {
	userID := GetUserID(req)
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}
}

// Notes model filter
func Notes(db *gorm.DB) *gorm.DB { // not working (why??)
	return db.Model(&models.Note{})
}

// NewFirst (by update)
func NewFirst(db *gorm.DB) *gorm.DB {
	return db.Order("updated_at DESC")
}
