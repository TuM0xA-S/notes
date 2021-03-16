package controllers

import (
	"contacts/auth"
	"contacts/models"
	"contacts/util"
	"encoding/json"
	"net/http"
)

//CreateAccount controller
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	a := &models.Account{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
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
	}

	resp := models.Login(a.Email, a.Password)

	util.Respond(w, resp)
}

//CreateContact for user controller
var CreateContact = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	contact := &models.Contact{}
	if err := json.NewDecoder(r.Body).Decode(contact); err != nil {
		resp := util.Message(false, "Invalid request")
		util.Respond(w, resp)
	}

	contact.UserID = user
	resp := contact.Create()
	util.Respond(w, resp)
})

//GetContacts for user controller
var GetContacts = auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)

	contacts := models.GetContacts(user)
	resp := util.Message(true, "Contacts fetched")
	resp["contacts"] = contacts

	util.Respond(w, resp)
})
