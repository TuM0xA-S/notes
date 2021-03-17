package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"notes/auth"
	"notes/models"
	"notes/util"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// MessageFromError ....
func MessageFromError(err error) map[string]interface{} {
	res := map[string]interface{}{}
	res["message"] = "OK"
	res["status"] = true

	if err != nil {
		res["message"] = err.Error()
		res["status"] = false
	}

	return res
}

//CreateAccount controller
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	a := &models.Account{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		log.Println(err, a)
		util.Respond(w, util.Message(false, "Invalid request"))
		return
	}

	err := a.Create()
	util.Respond(w, MessageFromError(err))
}

// GetUserID from request
func GetUserID(req *http.Request) uint {
	return req.Context().Value(auth.UserID).(uint)
}

//Login controller
func Login(w http.ResponseWriter, r *http.Request) {
	a := &models.Account{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		util.Respond(w, util.Message(false, "Invalide request"))
		return
	}

	accessToken, err := a.Login()
	resp := MessageFromError(err)
	if err == nil {
		resp["access_token"] = accessToken
	}
	util.Respond(w, resp)
}

//CreateNote for user controller
var CreateNote = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	note := &models.Note{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(note); err != nil {
		resp := util.Message(false, "Invalid request")
		util.Respond(w, resp)
		return
	}

	note.UserID = GetUserID(r)
	err := note.Create()
	util.Respond(w, MessageFromError(err))
})

//GetNotes for user controller
var GetNotes = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {

	notes := models.GetNotes(GetUserID(r))
	resp := util.Message(true, "OK")
	resp["notes"] = notes

	util.Respond(w, resp)
})

//NoteDetails ....
var NoteDetails = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	noteID, _ := strconv.Atoi(mux.Vars(r)["note_id"])

	note := &models.Note{}
	err := models.GetDB().First(note, "id = ? and user_id = ?", noteID, userID).Error
	if err == gorm.ErrRecordNotFound {
		util.Respond(w, util.Message(false, "no such note"))
	} else if err != nil {
		util.Respond(w, util.Message(false, "error with db"))
	} else {
		resp := util.Message(true, "OK")
		resp["note"] = note
		util.Respond(w, resp)
	}
})

//NoteRemove ....
var NoteRemove = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	noteID, _ := strconv.Atoi(mux.Vars(r)["note_id"])

	note := &models.Note{}
	err := models.GetDB().First(note, "id = ? and user_id = ?", noteID, userID).Error
	if err == gorm.ErrRecordNotFound {
		util.Respond(w, util.Message(false, "no such note"))
	} else if err != nil {
		util.Respond(w, util.Message(false, "error with db"))
	} else {
		models.GetDB().Delete(note)
		resp := util.Message(true, "OK")
		util.Respond(w, resp)
	}
})

//UserDetails ....
var UserDetails = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	user := &models.Account{}
	models.GetDB().Take(user, userID)
	resp := util.Message(true, "OK")
	resp["user"] = user
	util.Respond(w, resp)
})
