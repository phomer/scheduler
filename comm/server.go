package comm

import (
	"errors"
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/phomer/scheduler/accounts"
	"github.com/phomer/scheduler/jobs"
	"github.com/phomer/scheduler/log"
	"github.com/phomer/scheduler/sig"
)

var JSON = "application/json"

type Server struct {
	// Config parameters
	Host string
	Port string

	// Communication
	web *http.Server

	// Global Datastores
	Auth   *accounts.Authentication
	Sched  *jobs.Scheduled
	Active *jobs.Active
}

func NewServer() *Server {

	// TODO: Move to Config file
	host := "127.0.0.1"
	port := "8000"

	return &Server{
		Host: host,
		Port: port,

		web: NewHttpServer(host, port),

		Auth:   accounts.NewAuthentication(),
		Sched:  jobs.NewScheduled(),
		Active: jobs.NewActive(),
	}
}

// Set up a server to handle http traffic
func NewHttpServer(host string, port string) *http.Server {
	return &http.Server{
		Handler:      NewRouter(),
		Addr:         host + ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

// Allow global access to the top-level server configuration
var current *Server

// Catch an error if init is messed up
func Global() *Server {
	if current == nil {
		log.Fatal("Startup", errors.New("Missing Context"))
	}
	return current
}

// Catch an error if this is called twice
func set_server(server *Server) {
	if current != nil {
		log.Fatal("Instance", errors.New("Duplicate"))
	}
	current = server
}

// Start receiving requests
func (server *Server) Start() {
	fmt.Println("Starting HTTP for " + server.web.Addr)

	// Reload Datastores
	server.Auth.Reload()
	server.Sched.Reload()

	// Make this visible to the package, handers need access to shared config
	set_server(server)

	sig.Initialize()
	sig.Catch(syscall.SIGHUP, HandleSighup)

	err := server.web.ListenAndServe()
	if err != nil {
		fmt.Println("FATAL: Web Server Error: ", err)
		return
	}

	fmt.Println("Web Server Shutdown Nicely")
}

func HandleSigint() {
	// Shutdown cleanly
	fmt.Println("SIGINT Caught")
}

func HandleSighup() {
	// Reload Accounts and Jobs
	fmt.Println("Reloading Accounts")

	server := Global()

	// Reload accounts
	server.Auth.Reload()

	log.Dump("Accounts", server.Auth.Map)

	// TODO: SIGHUP seems to hang the web server, not sure if this is
	// a reasonable way to restart it, or it's just leaking ...
	// Probably need a better way to tell the server to reload
	err := server.web.ListenAndServe()
	if err != nil {
		fmt.Println("FATAL: Web Server Error: ", err)
		return
	}
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", Immediate)
	router.HandleFunc("/schedule", Schedule)
	router.HandleFunc("/tail", Tail)
	router.HandleFunc("/output", Output)
	router.HandleFunc("/status", Status)
	router.HandleFunc("/remove", Remove)

	return router
}
