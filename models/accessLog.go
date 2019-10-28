package models

import (
	"time"

	u "github.com/hrshadhin/license-server/utils"
)

// a struct to represent key access log
type KeyAccessLog struct {
	Domain    string       `gorm:"not null;index" json:"domain"`
	Key       string       `gorm:"not null" json:"key"`
	Referrer    string     `json:"referrer"`
	UserAgent   string     `json:"referrer"`
	RequestedAt time.Time  `gorm:"not null" json:"requested_at"`
	Status bool 		   `gorm:"not null" json:"update_key"`
	Message       string   `json:"message"`
}


func (keyAccessLog *KeyAccessLog) Create() bool {
	GetDB().Create(keyAccessLog)
	return true
}

func FetchAllKeyAccessLog() map[string]interface{} {
	var logs []KeyAccessLog
	err := GetDB().Table("key_access_log").Order("domain").Find(&logs).Error
	if err != nil {
		return u.Message(false, "Connection error. Please retry")
	}

	resp := u.Message(true, "Key access logs")
	resp["logs"] = logs
	return resp
}
