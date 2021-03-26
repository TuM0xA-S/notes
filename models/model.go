package models

import "time"

//Model is base for my models
type Model struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ErrValidation should be used to indicate a validation error
type ErrValidation string

func (e ErrValidation) Error() string {
	return string(e)
}

// IsErrValidation ...
func IsErrValidation(err error) bool {
	_, ok := err.(ErrValidation)
	return ok
}
