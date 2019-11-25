package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hrshadhin/license-server/models"
	u "github.com/hrshadhin/license-server/utils"
)

var CreateUser = func(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := user.Create() //Create account
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(user.Email, user.Password)
	u.Respond(w, resp)
}

var UserList = func(w http.ResponseWriter, r *http.Request) {

	resp := models.FetchAllUsers()
	u.Respond(w, resp)
}

var GetUser = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resp := models.GetUser(vars["userId"])
	u.Respond(w, resp)

}
var UpdateUser = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.UpdateUser(userId, user.Name, user.Email)
	u.Respond(w, resp)

}

var ChangePassword = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.ChangePassword(userId, user.Password)
	u.Respond(w, resp)

}

var DeleteUser = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resp := models.DeleteUser(vars["userId"])
	u.Respond(w, resp)

}
