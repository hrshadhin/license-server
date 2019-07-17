package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/hrshadhin/license-server/models"
	u "github.com/hrshadhin/license-server/utils"
)

var CreateKey = func(w http.ResponseWriter, r *http.Request) {

	key := &models.Key{}
	err := json.NewDecoder(r.Body).Decode(key) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := key.Create() //Create key
	u.Respond(w, resp)
}

var KeyList = func(w http.ResponseWriter, r *http.Request) {
	resp := models.FetchAllKeys()
	u.Respond(w, resp)
}
