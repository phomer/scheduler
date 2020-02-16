package comm

import (
	"errors"
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
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

	// Start the background service
	go jobs.ProcessSchedule()

	// Make this visible to the package, handers need access to shared config
	set_server(server)

	sig.Initialize()
	sig.Catch(syscall.SIGHUP, HandleSighup)

	// TODO: Quick fix
	go WatchAuthChange()

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

	// TODO: SIGHUP seems to hang the web server, not sure if this is
	// a reasonable way to restart it, or it's just leaking ...
	// Probably need a better way to tell the server to reload
	fmt.Println("Restarting Web Server after Signal")
	err := server.web.ListenAndServe()
	if err != nil {
		fmt.Println("FATAL: Web Server Error: ", err)
		return
	}
}

// TOOD: This is pretty horrible...
func WatchAuthChange() {
	server := Global()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("FileWatcher", err)
	}

	auth := server.Auth
	watcher.Add(auth.GetFilepath())
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Fatal("File Watcher Failed", errors.New("FileWatcher"))
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("Datamodified, reloading")
				auth.Reload()

				// Causes the web server to die
				err := server.web.ListenAndServe()
				if err != nil {
					log.Fatal("Web Server Error: ", err)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				log.Fatal("File Watcher Errors Failed", errors.New("FileWatcher"))
			}
			log.Fatal("File Watcher", err)
		}
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
