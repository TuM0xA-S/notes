package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"notes/auth"
	"notes/models"
	"notes/util"
)

//CreateAccount controller
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	a := &models.Account{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		log.Println(err, a)
		util.Respond(w, util.Message(false, "Invalid request"))
		return
	}

	resp := a.Create()
	util.Respond(w, resp)
}

//Login controller
func Login(w http.ResponseWriter, r *http.Request) {
	a := &models.Account{}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		util.Respond(w, util.Message(false, "Invalide request"))
		return
	}

	resp := models.Login(a.Username, a.Password)

	util.Respond(w, resp)
}

//CreateNote for user controller
var CreateNote = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	note := &models.Note{}
	if err := json.NewDecoder(r.Body).Decode(note); err != nil {
		resp := util.Message(false, "Invalid request")
		util.Respond(w, resp)
		return
	}

	note.UserID = user
	resp := note.Create()
	util.Respond(w, resp)
})

//GetNotes for user controller
var GetNotes = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)

	notes := models.GetNotes(user)
	resp := util.Message(true, "Notes fetched")
	resp["notes"] = notes

	util.Respond(w, resp)
})
