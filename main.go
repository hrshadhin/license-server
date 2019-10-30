package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hrshadhin/license-server/controllers"
	"github.com/hrshadhin/license-server/middleware"

	"github.com/getsentry/sentry-go"
	u "github.com/hrshadhin/license-server/utils"
	newrelic "github.com/newrelic/go-agent"
	nrgorilla "github.com/newrelic/go-agent/_integrations/nrgorilla/v1"
)

func main() {

	sentry.Init(sentry.ClientOptions{
		Dsn: u.MustGetEnv("sentry_dns"),
	})

	cfg := newrelic.NewConfig("License Server", u.MustGetEnv("new_relic_license_key"))
	app, err := newrelic.NewApplication(cfg)
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", controllers.Welcome).Methods("GET")
	router.HandleFunc("/api", controllers.Welcome).Methods("GET")
	router.HandleFunc("/api/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/users", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/api/users", controllers.UserList).Methods("GET")
	router.HandleFunc("/api/keys", controllers.KeyList).Methods("GET")
	router.HandleFunc("/api/keys", controllers.CreateKey).Methods("POST")
	router.HandleFunc("/api/keys/{domain}", controllers.UpdateKey).Methods("PATCH")
	router.HandleFunc("/api/verify", controllers.VerifyKey).Methods("POST")

	router.NotFoundHandler = http.HandlerFunc(middleware.NotFoundHandler)
	router.Use(middleware.JwtAuthentication) //attach JWT auth middleware

	port := os.Getenv("app_port")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println("Listening on http://localhost:" + port)

	err = http.ListenAndServe(":"+port, nrgorilla.InstrumentRoutes(router, app))
	if err != nil {
		fmt.Print(err)
	}
}
