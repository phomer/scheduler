package comm

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/phomer/scheduler/accounts"
	"github.com/phomer/scheduler/datastore"
	"github.com/phomer/scheduler/jobs"
	"github.com/phomer/scheduler/log"
)

// Execute a command right away
func Immediate(response http.ResponseWriter, request *http.Request) {
	// TODO: Hide this ...
	defer func() {
		err := recover()
		if err != nil {
			// TODO: Send a real error message
			response.Write([]byte("Server Error"))
		}
	}()

	account, err := ValidateRequest(response, request)
	if err != nil {
		log.Fatal("Validate", err)
	}

	cmd := GetClientRequest(request)

	jobs.Spawn(server(), account, cmd.Command, cmd.Args)

	// TODO: Stream file contents
	response.Write([]byte("Running Command"))
}

// Execute a command right away
func Schedule(response http.ResponseWriter, request *http.Request) {
	// TODO: Hide this ...
	defer func() {
		err := recover()
		if err != nil {
			// TODO: Send a real error message
			response.Write([]byte("Server Error"))
		}
	}()

	account, err := ValidateRequest(response, request)
	if err != nil {
		log.Fatal("Validate", err)
	}

	cmd := GetClientRequest(request)

	// TODO: Add to the scheduled jobs
	_ = cmd
	_ = account
	// TODO: Return the jobid

	jobid := []byte("42")

	response.Write(jobid)
}

// Execute a command right away
func Remove(response http.ResponseWriter, request *http.Request) {
	// TODO: Hide this ...
	defer func() {
		err := recover()
		if err != nil {
			// TODO: Send a real error message
			response.Write([]byte("Server Error"))
		}
	}()

	account, err := ValidateRequest(response, request)
	if err != nil {
		log.Fatal("Validate", err)
	}

	cmd := GetClientRequest(request)

	// TODO: Find the scheduled job
	_ = cmd
	_ = account
	// TODO: Remove it

	jobid := []byte("42")

	response.Write(jobid)
}

// Execute a command right away
func Tail(response http.ResponseWriter, request *http.Request) {
	// TODO: Hide this ...
	defer func() {
		err := recover()
		if err != nil {
			// TODO: Send a real error message
			response.Write([]byte("Server Error"))
		}
	}()

	account, err := ValidateRequest(response, request)
	if err != nil {
		log.Fatal("Validate", err)
	}

	cmd := GetClientRequest(request)

	// TODO: Find the running job
	_ = cmd
	_ = account
	// TODO: stream the output file, starting 200 lines up

	// TODO: Stream file contents
	response.Write([]byte("Running Command"))
}

// Execute a command right away
func Output(response http.ResponseWriter, request *http.Request) {
	// TODO: Hide this ...
	defer func() {
		err := recover()
		if err != nil {
			// TODO: Send a real error message
			response.Write([]byte("Server Error"))
		}
	}()

	account, err := ValidateRequest(response, request)
	if err != nil {
		log.Fatal("Validate", err)
	}

	cmd := GetClientRequest(request)

	// TODO: Find the running job
	_ = cmd
	_ = account
	// TODO: stream the output file, starting from the top

	// TODO: Stream file contents
	response.Write([]byte("Running Command"))
}

// Execute a command right away
func Status(response http.ResponseWriter, request *http.Request) {
	// TODO: Hide this ...
	defer func() {
		err := recover()
		if err != nil {
			// TODO: Send a real error message
			response.Write([]byte("Server Error"))
		}
	}()

	account, err := ValidateRequest(response, request)
	if err != nil {
		log.Fatal("Validate", err)
	}

	cmd := GetClientRequest(request)

	// TODO: Find the running job, if still in memory, check fs otherwise
	_ = cmd
	_ = account
	// TODO: get the exit code

	exitCode := []byte("42")
	response.Write(exitCode)
}

// Return an account that matches, or issue an error
func ValidateRequest(response http.ResponseWriter, request *http.Request) (*accounts.Account, error) {
	fmt.Println("Validating Request")

	username := accounts.Username() // TODO: Opps!

	isValid := true
	if !isValid {
		buffer := []byte("Authenication Failed")
		response.Write(buffer)
		return nil, errors.New("Authentication Failed")
	}

	root := server()
	account := root.Auth.Find(username)
	if account == nil {
		fmt.Println("Missing Account", username)
		log.Dump("Authentication", root.Auth)
		return nil, errors.New("User not Registered")
	}

	fmt.Println("Validated Account Access for ", account.Username)

	return account, nil
}

// Get the Request structure from HTTP request
func GetClientRequest(request *http.Request) *Request {

	buffer, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatal("Request Body", err)
	}

	log.Dump("Request Buffer", buffer)

	var prototype Request
	result := datastore.Deserialize(buffer, &prototype)

	return result.(*Request)
}
