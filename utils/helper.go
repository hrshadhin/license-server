package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func GetUserIpAddress(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	ipAndPortPart := strings.Split(IPAddress, ":")
	if len(ipAndPortPart) > 0 {
		IPAddress = ipAndPortPart[0]
	}
	return IPAddress
}

func GetUserAgent(r *http.Request) string {

	return r.Header.Get("User-Agent")
}