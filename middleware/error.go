package middleware

import (
	"net/http"

	u "github.com/hrshadhin/license-server/utils"
)

var NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	u.Respond(w, u.Message(false, "This resources was not found on our server"))
}
