package models

import (
	"log"
	"notes/util"

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
func (c *Note) Validate() (map[string]interface{}, bool) {
	if c.Title == "" {
		return util.Message(false, "Validation error. Title is empty"), false
	}
	if c.Body == "" {
		return util.Message(false, "Validation error. Body is empty"), false
	}

	if c.UserID <= 0 {
		return util.Message(false, "Validation error. UserID is invalid"), false
	}

	return util.Message(true, "Validation OK"), true
}

//Create note
func (c *Note) Create() map[string]interface{} {
	if resp, ok := c.Validate(); !ok {
		return resp
	}

	if GetDB().Create(c).Error != nil {
		return util.Message(false, "Failed to create")
	}

	resp := util.Message(true, "Succesfully created")
	resp["note"] = c

	return resp
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
