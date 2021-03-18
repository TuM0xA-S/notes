package models

import "time"

//Model is base for my models
type Model struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
}
