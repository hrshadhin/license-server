package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/hrshadhin/license-server/controllers"
	"github.com/hrshadhin/license-server/middleware"
	u "github.com/hrshadhin/license-server/utils"

	"github.com/getsentry/sentry-go"
	newrelic "github.com/newrelic/go-agent"
	nrgorilla "github.com/newrelic/go-agent/_integrations/nrgorilla/v1"
)

var (
	listenHost string
	listenPort string
	listenAddr string
)

func main() {

	flag.StringVar(&listenHost, "host", "127.0.0.1", "server listen host")
	flag.StringVar(&listenPort, "port", "8000", "server listen port")
	flag.Parse()

	listenAddr = listenHost + ":" + listenPort

	logger := log.New(os.Stdout, "LServer: ", log.LstdFlags)

	//make channels
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	// Catch the Ctrl-C and SIGTERM from kill command
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)

	//sentry logging
	sentry.Init(sentry.ClientOptions{
		Dsn: u.MustGetEnv("sentry_dns"),
	})

	//define new server
	server := newWebserver(logger)

	//run go routine for server shutdown signal
	go shutdownWebserver(server, logger, quit, done)

	//now start the server
	logger.Println("Server is ready to handle requests at", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Println("Server stopped")

}

func shutdownWebserver(server *http.Server, logger *log.Logger, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	logger.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(done)
}

func newWebserver(logger *log.Logger) *http.Server {
	//routes
	router := mux.NewRouter()
	router.HandleFunc("/", controllers.Welcome).Methods("GET")
	router.HandleFunc("/api", controllers.Welcome).Methods("GET")
	router.HandleFunc("/api/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/users", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/api/users", controllers.UserList).Methods("GET")
	router.HandleFunc("/api/users/{userId}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/api/users/{userId}", controllers.UpdateUser).Methods("PATCH")
	router.HandleFunc("/api/change-password/{userId}", controllers.ChangePassword).Methods("PATCH")
	router.HandleFunc("/api/users/{userId}", controllers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/api/keys", controllers.KeyList).Methods("GET")
	router.HandleFunc("/api/keys", controllers.CreateKey).Methods("POST")
	router.HandleFunc("/api/keys/{keyId}", controllers.GetKey).Methods("get")
	router.HandleFunc("/api/keys/{domain}", controllers.UpdateKey).Methods("PATCH")
	router.HandleFunc("/api/keys/{keyId}", controllers.DeleteKey).Methods("DELETE")
	router.HandleFunc("/api/verify", controllers.VerifyKey).Methods("POST")
	router.NotFoundHandler = http.HandlerFunc(middleware.NotFoundHandler)
	router.Use(middleware.JwtAuthentication) //attach JWT auth middleware

	//newrelic handler
	cfg := newrelic.NewConfig("License Server", u.MustGetEnv("new_relic_license_key"))
	app, err := newrelic.NewApplication(cfg)
	if nil != err {
		logger.Fatalf("Newrelic couldn't integrate: %v\n", err)
		os.Exit(1)
	}

	final_router := nrgorilla.InstrumentRoutes(router, app)

	//server instance
	return &http.Server{
		Addr:         listenAddr,
		Handler:      final_router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}
