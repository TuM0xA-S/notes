package models

import (
	"contacts/auth"
	"contacts/util"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//Account model
type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token" sql:"-"`
}

//Validate validates account data
func (a *Account) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(a.Email, "@") {
		return util.Message(false, "Email address is required"), false
	}

	if len(a.Password) < 6 {
		return util.Message(false, "Password is required"), false
	}

	temp := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", a.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return util.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return util.Message(false, "Connection error. Email address already in use by another user"), false
	}

	return util.Message(true, "Requitement passed"), true
}

//Create account in db
func (a *Account) Create() map[string]interface{} {
	if resp, ok := a.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	a.Password = string(hashedPassword)

	if GetDB().Create(a).Error != nil {
		return util.Message(false, "Connection error. Failed to create account")
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &auth.Token{UserID: a.ID})

	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	a.Token = tokenString

	a.Password = ""

	resp := util.Message(true, "Account has been created")
	resp["account"] = a

	return resp
}

//Login user
func Login(email, password string) map[string]interface{} {
	a := &Account{}
	if err := GetDB().Table("accounts").Where("email = ?", email).First(a).Error; err == gorm.ErrRecordNotFound {
		return util.Message(false, "Email address not found")
	} else if err != nil {
		return util.Message(false, "Connection error. Please retry")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)); err != nil {
		return util.Message(false, "Invalid login credentials. Please retry")
	}

	a.Password = ""

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &auth.Token{UserID: a.ID})
	a.Token, _ = token.SignedString([]byte(os.Getenv("token_password")))

	resp := util.Message(true, "Logged In")
	resp["account"] = a
	return resp
}

//GetUser by id
func GetUser(u uint) *Account {
	a := &Account{}
	if err := GetDB().Table("accounts").Where("id = ?", u).First(a).Error; err != nil {
		return nil
	}

	a.Password = ""
	return a
}
