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
	fmt.Println(temp.ExpiredAt)
	if(temp.UpdateKey){
		domainWithPad := key.Domain + fmt.Sprintf("%v", time.Now().Unix())
		hasher := sha1.New()
		hasher.Write([]byte(domainWithPad))
		newKey := hex.EncodeToString(hasher.Sum(nil))
		temp.Key = newKey
	}

	//if(temp.ExpiredAt == nil){
	//	temp.ExpiredAt = gorm.ex
	//}

	//update db
	models.GetDB().Model(&key).Updates(temp)

	u.Respond(w, resp)
	return

}