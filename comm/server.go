package comm

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	web *http.Server
}

func NewServer() *Server {
	return &Server{
		web: NewHttpServer(),
	}
}

func NewHttpServer() *http.Server {
	return &http.Server{
		Handler:      NewRouter(),
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func (server *Server) Start() {
	fmt.Println("Starting Web Server")

	err := server.web.ListenAndServe()
	if err != nil {
		fmt.Println("Web Server Error: ", err)
		return
	}

	fmt.Println("Web Server Shutdown Nicely")
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", Splat)

	return router
}

func Splat(response http.ResponseWriter, request *http.Request) {
	text := request.Body
	_ = text

	buffer := []byte("I'm alive...")

	response.Write(buffer)
}
