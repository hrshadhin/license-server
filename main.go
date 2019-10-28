package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hrshadhin/license-server/controllers"
	"github.com/hrshadhin/license-server/middleware"
)

func main() {

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

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
