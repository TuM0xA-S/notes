package models

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// GetDB returns db
func GetDB() *gorm.DB {
	return db
}

func init() {
	conn, err := gorm.Open(mysql.Open(os.Getenv("db_uri")))
	if err != nil {
		log.Fatalf("when connecting to db: %v", err)
	}
	db = conn
	db.Debug().AutoMigrate(&Account{}, &Contact{})
}
