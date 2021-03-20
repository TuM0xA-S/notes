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
	a := &models.User{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		util.RespondWithError(w, 400, "invalid request")
		return
	}

	err := a.Create()

	if models.IsErrValidation(err) {
		util.RespondWithError(w, 422, err.Error())
		return
	} else if err != nil {
		panic(err)
	}

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())

}

// GetUserID from request
func GetUserID(req *http.Request) uint {
	return req.Context().Value(auth.UserID).(uint)
}

//Login controller
func Login(w http.ResponseWriter, r *http.Request) {
	a := &models.User{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		util.RespondWithError(w, 400, "invalid request")
		return
	}

	accessToken, err := a.Login()
	if models.IsErrValidation(err) {
		util.RespondWithError(w, 422, err.Error())
		return
	} else if err != nil {
		panic(err)
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
	if models.IsErrValidation(err) {
		util.RespondWithError(w, 422, err.Error())
		return
	} else if err != nil {
		panic(err)
	}

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())
})

//NotesList for user controller
var NotesList = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	query := models.GetDB().Model(&models.Note{}).Select("id").Where("user_id = ?", GetUserID(r)).Order("updated_at DESC")
	page := GetPage(r)
	notes, err := Page(query, page)
	if err != nil {
		panic(err)
	}
	resp := util.ResponseBaseOK()
	resp["notes"] = notes
	resp["pagination"] = PaginationData(page, models.GetDB().Model(&models.Note{}))

	util.RespondWithJSON(w, 200, resp)
})

//PublishedNotesList ...
var PublishedNotesList = func(w http.ResponseWriter, r *http.Request) {
	query := models.GetDB().Model(&models.Note{}).Select("id").Where("published").Order("updated_at DESC")
	page := GetPage(r)
	notes, err := Page(query, page)
	if err != nil {
		panic(err)
	}
	resp := util.ResponseBaseOK()
	resp["notes"] = notes
	resp["pagination"] = PaginationData(page, models.GetDB().Model(&models.Note{}))

	util.RespondWithJSON(w, 200, resp)
}

// PublishedNoteDetail ...
func PublishedNoteDetail(w http.ResponseWriter, r *http.Request) {
	var noteID uint
	fmt.Sscan(mux.Vars(r)["note_id"], &noteID)
	note := &models.Note{Model: models.Model{ID: noteID}, Published: true}
	err := note.Get()
	if err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such note")
	} else if err != nil {
		panic(err)
	}
	resp := util.ResponseBaseOK()
	resp["note"] = note

	util.RespondWithJSON(w, 200, resp)
}

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
		panic(err)
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
	user := &models.User{}
	user.ID = userID
	if err := user.Get(); err != nil {
		panic("user should always be valid because of authorization")
	}
	user.Password = "<hashed>"
	resp := util.ResponseBaseOK()
	resp["user"] = user
	util.RespondWithJSON(w, 200, resp)
})
