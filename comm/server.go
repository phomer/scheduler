package comm

import (
	"errors"
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/phomer/scheduler/accounts"
	"github.com/phomer/scheduler/log"
	"github.com/phomer/scheduler/sig"
)

var JSON = "application/json"

type Server struct {
	Auth *accounts.Authentication
	web  *http.Server
	Host string
	Port string
}

func NewServer() *Server {

	// TODO: Move to Config file
	host := "127.0.0.1"
	port := "8000"

	return &Server{
		Auth: accounts.NewAuthentication(),
		web:  NewHttpServer(host, port),
		Host: host,
		Port: port,
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
func server() *Server {
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
	fmt.Println("Starting HTTP")

	// TODO: Needs cleanup ...
	// Reload Accounts
	server.Auth.Lock()
	server.Auth.Reload()
	server.Auth.Unlock()

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

	server := server()

	server.Auth.Lock()
	server.Auth.Reload()
	server.Auth.Unlock()

	log.Dump("Accounts", server.Auth.Map)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", Immediate)

	return router
}
