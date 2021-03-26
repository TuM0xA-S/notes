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

// NotFound Handler ..
func NotFound(w http.ResponseWriter, r *http.Request) {
	util.RespondWithError(w, 404, r.URL.String()+" not found")
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
	notes := &[]map[string]interface{}{}
	err := models.GetDB().Model(&models.Note{}).Scopes(OwnedBy(r), Paginate(r), NewFirst).
		Omit("body").Find(notes).Error
	if err != nil {
		panic(err)
	}
	resp := util.ResponseBaseOK()
	resp["notes"] = notes
	resp["pagination"] = PaginationData(r, models.GetDB().Model(&models.Note{}).Scopes(OwnedBy(r)))

	util.RespondWithJSON(w, 200, resp)
})

//PublishedNotesList ...
var PublishedNotesList = func(w http.ResponseWriter, r *http.Request) {
	notes := &[]map[string]interface{}{}
	err := models.GetDB().Model(&models.Note{}).Scopes(Published, Paginate(r), NewFirst).
		Omit("body").Find(notes).Error
	if err != nil {
		panic(err)
	}
	resp := util.ResponseBaseOK()
	resp["notes"] = notes
	resp["pagination"] = PaginationData(r, models.GetDB().Model(&models.Note{}).Scopes(Published))

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
	resp["note"] = notePrework(r, note)

	util.RespondWithJSON(w, 200, resp)
}

// AnotherUserDetail ...
func AnotherUserDetail(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	fmt.Sscan(mux.Vars(r)["user_id"], &user.ID)

	err := user.Get()
	if err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such user")
	} else if err != nil {
		panic(err)
	}
	user.Password = ""
	resp := util.ResponseBaseOK()
	resp["user"] = user
	util.RespondWithJSON(w, 200, resp)
}

func notePrework(r *http.Request, n *models.Note) *models.Note {
	if r.URL.Query().Get("no_body") == "true" {
		n.Body = ""
	}
	return n
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
		resp["note"] = notePrework(r, note)
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
		panic(err)
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
	patch := &models.NotePatch{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(patch); err != nil {
		util.RespondWithError(w, 400, "invalid request")
		return
	}

	err := note.Update(patch)

	if err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such note")
	} else if models.IsErrValidation(err) {
		util.RespondWithError(w, 422, err.Error())
	} else if err != nil {
		panic(err)
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
