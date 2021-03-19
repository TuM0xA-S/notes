package models

import (
	"fmt"
	"notes/auth"
	. "notes/config"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//User model
type User struct {
	Model
	Username string `json:"username"`
	Password string `json:"password"`
	Notes    []Note `json:"-"`
}

//Validate validates account data
func (a *User) Validate() error {
	if len(a.Username) < 4 || len(a.Username) > 20 {
		return fmt.Errorf("Username is required(4 <= len <= 20)")
	}

	if len(a.Password) < 6 || len(a.Password) > 30 {
		return fmt.Errorf("Password is required(6 <= len <= 30)")
	}

	temp := &User{}
	err := GetDB().Where("username = ?", a.Username).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("Connection error. Please retry")
	}
	if temp.Username != "" {
		return fmt.Errorf("Username is already in use")
	}
	return nil
}

// HashPassword generates hash for password (WOW)
func HashPassword(password string) string {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hash := string(hashBytes)
	return hash
}

//Create account in db
func (a *User) Create() error {
	if err := a.Validate(); err != nil {
		return err
	}

	a.Password = HashPassword(a.Password)

	if GetDB().Create(a).Error != nil {
		return fmt.Errorf("Connection error. Failed to create account")
	}

	return nil
}

// GenerateToken for user
func GenerateToken(uid uint) string {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &auth.Token{UserID: uid})
	accessToken, _ := token.SignedString([]byte(Cfg.TokenPassword))
	return accessToken
}

//Login user
func (a *User) Login() (string, error) {
	password := a.Password
	if err := GetDB().Where("username = ?", a.Username).First(a).Error; err == gorm.ErrRecordNotFound {
		return "", fmt.Errorf("Username not found")
	} else if err != nil {
		return "", fmt.Errorf("Connection error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("Invalid login credentials")
	}

	return GenerateToken(a.ID), nil
}

// Get user
func (a *User) Get() error {
	return GetDB().Take(a, a.ID).Error
}
