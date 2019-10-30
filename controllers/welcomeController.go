package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	u "github.com/hrshadhin/license-server/utils"
	"github.com/newrelic/go-agent"
	"github.com/getsentry/sentry-go"
)

var Welcome = func(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path
	config := newrelic.NewConfig("License Server", "182f058bf3147cf1cf16733cf834a1ff0e547433")
	app, err := newrelic.NewApplication(config)
	if err != nil {
		message := "newrelic not bootup!"
		sentry.CaptureException(errors.New(message))
		sentry.Flush(time.Second * 1)
		fmt.Println(message)
	}
	fmt.Println(requestPath)
	txn := app.StartTransaction(requestPath, w, r)
	defer txn.End()
	u.Respond(w, u.Message(true, "Welcome to license server api!"))
}
