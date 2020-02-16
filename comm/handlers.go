package comm

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/phomer/scheduler/accounts"
	"github.com/phomer/scheduler/datastore"
	"github.com/phomer/scheduler/jobs"
	"github.com/phomer/scheduler/log"
)

// Execute a command right away
func Immediate(response http.ResponseWriter, request *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			HandleError(response, err)
		}
	}()

	account, err := ValidateRequest(response, request)
	if err != nil {
		HandleError(response, err)
	}

	cmd := GetClientRequest(request)

	command := jobs.NewCommand(cmd)

	job := jobs.NewImmediateJob(account, command)

	jobs.Spawn(account, job)

	StreamResponse(job, response)
}

// Execute a command right away
func Schedule(response http.ResponseWriter, request *http.Request) {
	// TODO: Hide this ...
	defer func() {
		err := recover()
		if err != nil {
			HandleError(response, err)
		}
	}()

	account, err := ValidateRequest(response, request)
	if err != nil {
		HandleError(response, err)
		return
	}

	cmd := GetClientRequest(request)

	command := jobs.NewCommand(cmd)

	job := jobs.NewActiveJob(account, command)

	response.Write(Int2Bytes(job.Cmd.JobId))
}

// Execute a command right away
func Remove(response http.ResponseWriter, request *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			HandleError(response, err)
		}
	}()

	account, err := ValidateRequest(response, request)
	if err != nil {
		HandleError(response, err)
		return
	}

	cmd := GetClientRequest(request)

	status := Global().Sched.RemoveUserCommand(account.Username, cmd.JobId)

	response.Write(Int2Bytes(status))
}

// Execute a command right away
func Tail(response http.ResponseWriter, request *http.Request) {
	// TODO: Hide this ...
	defer func() {
		err := recover()
		if err != nil {
			HandleError(response, err)
		}
	}()

	account, err := ValidateRequest(response, request)
	if err != nil {
		HandleError(response, err)
		return
	}

	cmd := GetClientRequest(request)

	entry := Global().Active.FindJobStatus(account.Username, cmd.JobId)

	StreamResponse(entry, response)
}

// Execute a command right away
func Output(response http.ResponseWriter, request *http.Request) {
	// TODO: Hide this ...
	defer func() {
		err := recover()
		if err != nil {
			HandleError(response, err)
		}
	}()

	account, err := ValidateRequest(response, request)
	if err != nil {
		HandleError(response, err)
		return
	}

	cmd := GetClientRequest(request)

	entry := Global().Active.FindJobStatus(account.Username, cmd.JobId)

	StreamResponse(entry, response)
}

// Execute a command right away
func Status(response http.ResponseWriter, request *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			HandleError(response, err)
		}
	}()

	account, err := ValidateRequest(response, request)
	if err != nil {
		HandleError(response, err)
		return
	}

	cmd := GetClientRequest(request)

	entry := Global().Active.FindJobStatus(account.Username, cmd.JobId)
	if entry != nil {
		if entry.IsRunning {
			response.Write([]byte("Running"))
		} else {
			response.Write([]byte("Finished"))
		}
	} else {
		response.Write([]byte("Finished"))
	}
}

// Handle logging and sending back and error
func HandleError(response http.ResponseWriter, err interface{}) {
	fmt.Println("Internal Error " + err.(error).Error())
	fmt.Println(string(debug.Stack()))

	response.Write([]byte("Server Error " + err.(error).Error()))
}

// Return an account that matches, or issue an error
func ValidateRequest(response http.ResponseWriter, request *http.Request) (*accounts.Account, error) {
	fmt.Println("Validating Request")

	tokenString := request.Header.Get("Authorization")
	if tokenString == "" {
		buffer := []byte("Missing Token, Authenication Failed")
		response.Write(buffer)
		return nil, errors.New("Missing Token, Authentication Failed")
	}

	token := accounts.NewToken(tokenString)

	log.Dump("Token", token)

	if !accounts.Validate(token) {
		buffer := []byte("Invalid Token, Authenication Failed")
		response.Write(buffer)
		return nil, errors.New("Invalid Token, Authentication Failed")
	}

	username := request.Header["Name"][0]
	if username == "" {
		fmt.Println("Missing Account", username)
		log.Dump("Authentication", Global().Auth)

		return nil, errors.New("User not Registered")
	}

	account := Global().Auth.Find(username)

	if account == nil {
		fmt.Println("Missing Account Information ", username)
		log.Dump("Authentication", Global().Auth)

		return nil, errors.New("User not Registered")
	}

	fmt.Println("Validated Account Access for ", account.Username)

	return account, nil
}

// Get the Request structure from HTTP request
func GetClientRequest(request *http.Request) *jobs.Request {

	buffer, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatal("Request Body", err)
	}

	log.Dump("Request Buffer", buffer)

	var prototype jobs.Request
	result := datastore.Deserialize(buffer, &prototype)

	return result.(*jobs.Request)
}

func StreamResponse(job *jobs.ActiveJob, response http.ResponseWriter) {
	// Stream what is there now
	// Watch and stream the rest
	var err error
	if job.File == nil {
		job.File, err = os.Open(job.Cmd.Filepath)
		if err != nil {
			HandleError(response, err)
			return
		}

		// TODO: Pick a good offset if the incomming option was tail
		job.Offset = 0
	}

	// Move us format to a good starting place
	if job.Offset != 0 {
		job.Offset, err = job.File.Seek(job.Offset, 0)
		if err != nil {
			HandleError(response, err)
		}
	}

	// Keep trying until the job stops, it runs out of file, or the pipe closes.
	for {
		buffer := make([]byte, 2048)

		count, err := job.File.Read(buffer)
		if err != nil && err != io.EOF {
			HandleError(response, err)
			return
		}

		if count == 0 {
			// TODO: We should switch to a nicer way of handling this, be we won't
			if !Global().Active.IsActive(job.Pid) {
				// It's still running, so it might still produce outout
				return

			} else {
				// TODO: Polling is ugly
				time.Sleep(1 * time.Second)
			}
		} else {
			response.Write(buffer)
			job.Offset += int64(count)
		}
	}
}

func Int2Bytes(value int) []byte {
	return []byte(fmt.Sprintf("%d", value))
}
