package models

import (
	"contacts/util"
	"log"

	"gorm.io/gorm"
)

//Contact with name, phone and owner
type Contact struct {
	gorm.Model
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	UserID uint   `json:"user_id"`
}

//Validate contact
func (c *Contact) Validate() (map[string]interface{}, bool) {
	if c.Name == "" {
		return util.Message(false, "Validation error. Name is empty"), false
	}
	if c.Phone == "" {
		return util.Message(false, "Validation error. Phone is empty"), false
	}

	if c.UserID <= 0 {
		return util.Message(false, "Validation error. UserID is invalid"), false
	}

	return util.Message(true, "Validation OK"), true
}

//Create contact
func (c *Contact) Create() map[string]interface{} {
	if resp, ok := c.Validate(); !ok {
		return resp
	}

	if GetDB().Create(c).Error != nil {
		return util.Message(false, "Failed to create")
	}

	resp := util.Message(true, "Succesfully created")
	resp["contact"] = c

	return resp
}

//GetContact by id
func GetContact(id uint) *Contact {
	c := &Contact{}
	err := GetDB().Table("contacts").Where("id = ?", id).First(c).Error
	if err != nil {
		return nil
	}

	return c
}

//GetContacts for user
func GetContacts(user uint) []*Contact {
	contacts := []*Contact{}
	err := GetDB().Table("contacts").Where("user_id = ?", user).Find(&contacts).Error
	log.Println(err)
	if err != nil {
		return nil
	}

	return contacts
}
