package models

import (
	"notes/auth"
	. "notes/config"
	"notes/util"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//Account model
type Account struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

//Validate validates account data
func (a *Account) Validate() (map[string]interface{}, bool) {
	if len(a.Username) < 4 {
		return util.Message(false, "Username is required(min len 4)"), false
	}

	if len(a.Password) < 6 {
		return util.Message(false, "Password is required(min len 6)"), false
	}

	temp := &Account{}
	err := GetDB().Table("accounts").Where("username = ?", a.Username).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return util.Message(false, "Connection error. Please retry"), false
	}
	if temp.Username != "" {
		return util.Message(false, "Username is already in use"), false
	}

	return util.Message(true, "Requitement passed"), true
}

// HashPassword generates hash for password (WOW)
func HashPassword(password string) string {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hash := string(hashBytes)
	return hash
}

//Create account in db
func (a *Account) Create() map[string]interface{} {
	if resp, ok := a.Validate(); !ok {
		return resp
	}

	a.Password = HashPassword(a.Password)

	if GetDB().Create(a).Error != nil {
		return util.Message(false, "Connection error. Failed to create account")
	}

	resp := util.Message(true, "Account has been created")
	resp["account"] = a

	return resp
}

func GenerateToken(uid uint) string {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &auth.Token{UserID: uid})
	accessToken, _ := token.SignedString([]byte(Cfg.TokenPassword))
	return accessToken
}

//Login user
func Login(username, password string) map[string]interface{} {
	a := &Account{}
	if err := GetDB().Table("accounts").Where("username = ?", username).First(a).Error; err == gorm.ErrRecordNotFound {
		return util.Message(false, "Username not found")
	} else if err != nil {
		return util.Message(false, "Connection error. Please retry")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)); err != nil {
		return util.Message(false, "Invalid login credentials. Please retry")
	}

	resp := util.Message(true, "Logged In")
	resp["access_token"] = GenerateToken(a.ID)
	a.Password = password
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
