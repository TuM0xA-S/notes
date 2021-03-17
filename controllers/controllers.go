package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"notes/auth"
	"notes/models"
	"notes/util"
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
	user := r.Context().Value("user").(uint)
	note := &models.Note{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(note); err != nil {
		resp := util.Message(false, "Invalid request")
		util.Respond(w, resp)
		return
	}

	note.UserID = user
	err := note.Create()
	util.Respond(w, MessageFromError(err))
})

//GetNotes for user controller
var GetNotes = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)

	notes := models.GetNotes(user)
	resp := util.Message(true, "OK")
	resp["notes"] = notes

	util.Respond(w, resp)
})
