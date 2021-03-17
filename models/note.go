package models

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

//Note with name, phone and owner
type Note struct {
	gorm.Model
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID uint   `json:"user_id"`
}

//Validate note
func (c *Note) Validate() error {
	if c.Title == "" {
		return fmt.Errorf("Validation error. Title is empty")
	}
	if c.Body == "" {
		return fmt.Errorf("Validation error. Body is empty")
	}

	if c.UserID <= 0 {
		return fmt.Errorf("Validation error. UserID is invalid")
	}

	return nil
}

//Create note
func (c *Note) Create() error {
	if err := c.Validate(); err != nil {
		return err
	}

	if GetDB().Create(c).Error != nil {
		return fmt.Errorf("Failed to create")
	}

	return nil
}

//GetNote by id
func GetNote(id uint) *Note {
	c := &Note{}
	err := GetDB().Table("notes").Where("id = ?", id).First(c).Error
	if err != nil {
		return nil
	}

	return c
}

//GetNotes for user
func GetNotes(user uint) []*Note {
	notes := []*Note{}
	err := GetDB().Table("notes").Where("user_id = ?", user).Find(&notes).Error
	log.Println(err)
	if err != nil {
		return nil
	}

	return notes
}
