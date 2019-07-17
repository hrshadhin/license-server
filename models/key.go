package models

import (
	"strings"
	"time"

	u "github.com/hrshadhin/license-server/utils"
	"github.com/jinzhu/gorm"
)

// a struct to represent key
type Key struct {
	ID        uint       `gorm:"primary_key"`
	Domain    string     `gorm:"not null;unique_index" json:"domain"`
	Key       string     `gorm:"not null" json:"key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ExpiredAt *time.Time `json:"expired_at"`
}

//Validate incoming user details...
func (key *Key) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(key.Domain, ".") {
		return u.Message(false, "Invalid domain"), false
	}

	//domain must be unique
	temp := &Key{}

	//check for errors and duplicate emails
	err := GetDB().Table("keys").Where("domain = ?", key.Domain).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Domain != "" {
		return u.Message(false, "Domain already exists!"), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (key *Key) Create() map[string]interface{} {

	if resp, ok := key.Validate(); !ok {
		return resp
	}

	hashStr := "This is super key"
	key.Key = hashStr

	GetDB().Create(key)

	if key.ID <= 0 {
		return u.Message(false, "Failed to create key, connection error.")
	}

	response := u.Message(true, "Key has been created")
	response["key"] = key
	return response
}

func FetchAllKeys() map[string]interface{} {

	var keys []Key
	err := GetDB().Table("keys").Order("id").Find(&keys).Error
	if err != nil {
		return u.Message(false, "Connection error. Please retry")
	}

	resp := u.Message(true, "Key List")
	resp["keys"] = keys
	return resp
}
