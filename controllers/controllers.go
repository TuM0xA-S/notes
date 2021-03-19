package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notes/auth"
	"notes/models"
	"notes/util"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

//CreateAccount controller
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	a := &models.Account{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		util.RespondWithError(w, 400, "invalid request")
		return
	}

	err := a.Create()

	if err != nil {
		util.RespondWithError(w, 422, err.Error())
		return
	}

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())

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
		util.RespondWithError(w, 400, "invalid request")
		return
	}

	accessToken, err := a.Login()
	if err != nil {
		util.RespondWithError(w, 422, err.Error())
		return
	}

	resp := util.ResponseBaseOK()
	resp["access_token"] = accessToken
	util.RespondWithJSON(w, 200, resp)
}

//CreateNote for user controller
var CreateNote = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	note := &models.Note{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(note); err != nil {
		util.RespondWithError(w, 400, "body cannot be used for create")
		return
	}

	note.UserID = GetUserID(r)
	err := note.Create()

	// idk how to manage error properly here
	if err != nil {
		util.RespondWithError(w, 400, err.Error())
		return
	}

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())
})

//GetNotes for user controller
var GetNotes = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	notes := &[]map[string]interface{}{}
	err := models.GetDB().Model(&models.Note{}).Select("title", "id").Find(notes, "user_id = ?", GetUserID(r)).Error
	if err != nil {
		panic("error with db")
	}
	resp := util.ResponseBaseOK()
	resp["notes"] = notes

	util.RespondWithJSON(w, 200, resp)
})

//NoteDetails ....
var NoteDetails = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	var noteID uint
	fmt.Sscan(mux.Vars(r)["note_id"], &noteID)

	note := &models.Note{}
	note.ID = noteID
	note.UserID = userID

	err := note.Get()

	if err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such note")
	} else if err != nil {
		panic("troubles with db")
	} else {
		resp := util.ResponseBaseOK()
		resp["note"] = note
		util.RespondWithJSON(w, 200, resp)
	}
})

//NoteRemove ....
var NoteRemove = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	var noteID uint
	fmt.Sscan(mux.Vars(r)["note_id"], &noteID)

	note := &models.Note{}
	note.ID = noteID
	note.UserID = userID

	err := note.Remove()

	if err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such note")
	} else if err != nil {
		panic("troubles with db")
	} else {
		util.RespondWithJSON(w, 200, util.ResponseBaseOK())
	}
})

//NoteUpdate ....
var NoteUpdate = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	var noteID uint
	fmt.Sscan(mux.Vars(r)["note_id"], &noteID)

	note := &models.Note{}
	note.ID = noteID
	note.UserID = userID
	patch := &models.Note{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(patch); err != nil {
		util.RespondWithError(w, 400, "invalid request")
		return
	}

	err := note.Update(patch)

	if err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such note")
	} else if err != nil {
		panic("troubles with db")
	} else {
		util.RespondWithJSON(w, 200, util.ResponseBaseOK())
	}
})

//UserDetails ....
var UserDetails = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	user := &models.Account{}
	user.ID = userID
	if err := user.Get(); err != nil {
		panic("user should always be valid because of authorization")
	}
	user.Password = "<hashed>"
	resp := util.ResponseBaseOK()
	resp["user"] = user
	util.RespondWithJSON(w, 200, resp)
})
