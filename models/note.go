package models

import (
	"fmt"
)

//Note with title, body and ownwer
type Note struct {
	Model
	Title     string `json:"title"`
	Body      string `json:"body"`
	UserID    uint   `json:"user_id"`
	Published bool   `json:"published"`
}

//Validate note
func (n *Note) Validate() error {
	if len(n.Title) < 3 || len(n.Title) > 40 {
		return ErrValidation("Validation error. Title len should be (3 <= len <= 40)")
	}
	if len(n.Body) > 512 {
		return ErrValidation("Validation error. Body is too big(max len 512)")
	}

	if n.UserID <= 0 {
		return ErrValidation("Validation error. UserID is invalid")
	}

	return nil
}

//Create note
func (n *Note) Create() error {
	if err := n.Validate(); err != nil {
		return err
	}

	if err := GetDB().Create(n).Error; err != nil {
		panic(fmt.Errorf("when creating in db: %v", err))
	}

	return nil
}

//Get note
func (n *Note) Get() error {
	return GetDB().Where(n).Take(n).Error
}

//Save note
func (n *Note) Save() error {
	return GetDB().Save(n).Error
}

//Remove note
func (n *Note) Remove() error {
	return GetDB().Where(n).Delete(n).Error
}

//Update note
func (n *Note) Update(patch *Note) error {
	if err := n.Get(); err != nil {
		return err
	}
	if patch.Body != "" {
		n.Body = patch.Body
	}
	if patch.Title != "" {
		n.Title = patch.Title
	}

	if err := n.Validate(); err != nil {
		return err
	}
	return GetDB().Save(n).Error
}
