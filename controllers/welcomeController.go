package controllers

import (
	"net/http"

	u "github.com/hrshadhin/license-server/utils"
)

var Welcome = func(w http.ResponseWriter, r *http.Request) {
	u.Respond(w, u.Message(true, "Welcome to license server api!"))
}
