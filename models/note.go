package models

import (
	"fmt"

	"gorm.io/gorm"
)

//Note with title, body and ownwer
type Note struct {
	gorm.Model
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID uint   `json:"user_id"`
}

//Validate note
func (n *Note) Validate() error {
	if len(n.Title) < 3 || len(n.Title) > 40 {
		return fmt.Errorf("Validation error. Title len should be (3 <= len <= 40)")
	}
	if len(n.Body) > 512 {
		return fmt.Errorf("Validation error. Body is too big(max len 512)")
	}

	if n.UserID <= 0 {
		return fmt.Errorf("Validation error. UserID is invalid")
	}

	return nil
}

//Create note
func (n *Note) Create() error {
	if err := n.Validate(); err != nil {
		return err
	}

	if GetDB().Create(n).Error != nil {
		return fmt.Errorf("Failed to create")
	}

	return nil
}

//GetNotes for user
func GetNotes(user uint) []*Note {
	notes := []*Note{}
	err := GetDB().Where("user_id = ?", user).Find(&notes).Error
	if err != nil {
		return nil
	}

	return notes
}
