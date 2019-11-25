package controllers

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
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

var UpdateKey = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resp := u.Message(true, "Key updated!")

	if !strings.Contains(vars["domain"], ".") {
		resp = u.Message(false, "Invalid domain")
	}

	key := &models.Key{}
	found := key.FindByDomain(vars["domain"])
	if !found {
		resp = u.Message(false, "Domain not found!")
	}

	temp := &models.Key{}
	err := json.NewDecoder(r.Body).Decode(temp)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	//just be safe about user will mass
	temp.Key = key.Key
	//fmt.Println(temp.ExpiredAt)
	if temp.UpdateKey {
		domainWithPad := key.Domain + fmt.Sprintf("%v", time.Now().Unix())
		hasher := sha1.New()
		hasher.Write([]byte(domainWithPad))
		key.Key = hex.EncodeToString(hasher.Sum(nil))
	}

	if temp.ExpiredAt != nil {
		key.ExpiredAt = temp.ExpiredAt
	}

	//update db
	models.GetDB().Model(&key).Updates(key)

	u.Respond(w, resp)
	return

}

var VerifyKey = func(w http.ResponseWriter, r *http.Request) {
	status := true
	message := "Verified"

	type DomainKey struct {
		Domain string
		Key    string
	}
	requestBody := &DomainKey{}
	err := json.NewDecoder(r.Body).Decode(requestBody)
	if err != nil {
		status = false
		message = "Invalid request!"
	}

	key := &models.Key{}
	found := key.FindByDomain(requestBody.Domain)
	if !found {
		status = false
		message = "Domain not registered!"
	}

	if key.Key != requestBody.Key {
		status = false
		message = "Invalid key!"
	}

	if key.ExpiredAt != nil {
		nowStamp := time.Now()
		expiredAt := *key.ExpiredAt
		if nowStamp.After(expiredAt) {
			status = false
			message = "License expired!"
		}
	}

	// log the event
	log := &models.KeyAccessLog{}
	log.Domain = requestBody.Domain
	log.Key = requestBody.Key
	log.RequestedAt = time.Now()
	log.Referrer = u.GetUserIpAddress(r)
	log.UserAgent = u.GetUserAgent(r)
	log.Status = status
	log.Message = message

	log.Create()

	u.Respond(w, u.Message(status, message))
	return

}

var GetKey = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resp := models.GetKey(vars["keyId"])
	u.Respond(w, resp)

}

var DeleteKey = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resp := models.DeleteKey(vars["keyId"])
	u.Respond(w, resp)

}
