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

//Validate validates account data(can it be created?)
func (a *User) Validate() error {
	if len(a.Username) < 4 || len(a.Username) > 20 {
		return ErrValidation("Username is required(4 <= len <= 20)")
	}

	if len(a.Password) < 6 || len(a.Password) > 30 {
		return ErrValidation("Password is required(6 <= len <= 30)")
	}

	err := GetDB().Where("username = ?", a.Username).First(&User{}).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return nil
	case nil:
		return ErrValidation("Username is already in use")
	default:
		panic(fmt.Errorf("when fetching from db: %v", err))
	}
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

	if err := GetDB().Create(a).Error; err != nil {
		panic(fmt.Errorf("when creating in db: %v", err))
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
	err := GetDB().Where("username = ?", a.Username).First(a).Error
	if err == gorm.ErrRecordNotFound {
		return "", ErrValidation("Username not found")
	} else if err != nil {
		panic(fmt.Errorf("when fetching from db: %v", err))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)); err != nil {
		return "", ErrValidation("Invalid login credentials")
	}

	return GenerateToken(a.ID), nil
}

// Get user
func (a *User) Get() error {
	return GetDB().Take(a, a.ID).Error
}
