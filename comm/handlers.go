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

	fmt.Println("Immediate Command")

	cmd := GetClientRequest(request)
	command := jobs.NewCommand(cmd)

	job := jobs.NewImmediateJob(account, command)
	err = jobs.Spawn(account, job)

	if err != nil {
		HandleError(response, err)
	} else {
		StreamResponse(job, response)
	}
}

// Execute a command right away
func Schedule(response http.ResponseWriter, request *http.Request) {

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

	command = Global().Sched.AddScheduledCommand(account.Username, command)
	fmt.Println("Scheduling Work", command.JobId)

	WriteResponse(response, "JobId: ", command.JobId)
}

func WriteResponse(response http.ResponseWriter, args ...interface{}) {
	buffer := fmt.Sprint(args...)
	response.Write([]byte(buffer))
	response.Write([]byte("\n"))
	response.(http.Flusher).Flush()
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

	WriteResponse(response, "Status: ", status)
}

// Execute a command right away
func Tail(response http.ResponseWriter, request *http.Request) {

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
		StreamResponse(entry, response)
	} else {
		WriteResponse(response, "Unknown Job Id: ", cmd.JobId)
	}
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
			WriteResponse(response, "Status: ", "Running")
		} else {
			WriteResponse(response, "Status: ", "Finished")
		}
	} else {
		WriteResponse(response, "Status: ", "Pending")
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

	tokenString := request.Header.Get("Authorization")
	if tokenString == "" {
		buffer := []byte("Missing Token, Authenication Failed")
		response.Write(buffer)
		return nil, errors.New("Missing Token, Authentication Failed")
	}

	token := accounts.NewToken(tokenString)

	if !accounts.Validate(token) {
		buffer := []byte("Invalid Token, Authenication Failed")
		response.Write(buffer)
		return nil, errors.New("Invalid Token, Authentication Failed")
	}

	username := request.Header["Name"][0]
	if username == "" {
		fmt.Println("Missing Account", username)

		return nil, errors.New("User not Registered")
	}

	account := Global().Auth.Find(username)

	if account == nil {
		fmt.Println("Missing Account Information ", username)

		return nil, errors.New("User not Registered")
	}

	return account, nil
}

// Get the Request structure from HTTP request
func GetClientRequest(request *http.Request) *jobs.Request {

	buffer, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatal("Request Body", err)
	}

	var prototype jobs.Request
	result := datastore.Deserialize(buffer, &prototype)

	return result.(*jobs.Request)
}

func StreamResponse(job *jobs.ActiveJob, response http.ResponseWriter) {
	if job == nil {
		WriteResponse(response, "Inactive Job")
		return
	}

	// Stream what is there now
	// Watch and stream the rest
	var err error
	if job.File == nil {
		fmt.Println("Opening File", job.Cmd.Filepath)

		job.File, err = os.Open(job.Cmd.Filepath)
		if err != nil {
			HandleError(response, err)
			return
		}

		// TODO: Pick a good offset if the incomming option was tail
		job.Offset = 0
	}

	// Move the start to a later location in the file.
	if job.Offset != 0 {
		fmt.Println("Moving Forward")
		job.Offset, err = job.File.Seek(job.Offset, 0)
		if err != nil {
			HandleError(response, err)
		}
	}

	// Keep trying until the job stops, it runs out of file, or the pipe closes.
	finished := false
	for {
		time.Sleep(1000 * time.Millisecond)
		buffer := make([]byte, 1024)

		count, err := job.File.Read(buffer)
		if err != nil {
			if err == io.EOF {
				if finished {
					// We're tried a couple of times and the process says
					// it is done.
					response.(http.Flusher).Flush()
					return
				}
				if !Global().Active.IsActive(job.Pid) {
					finished = true
				}
			} else {
				HandleError(response, err)
				return
			}
		} else {
			finished = false
		}

		if count != 0 {
			response.Write(buffer)
			response.(http.Flusher).Flush()
			job.Offset += int64(count)
		} else {
			// Nothing to read write now
			fmt.Println("Waiting")
			time.Sleep(2 * time.Millisecond)
		}
	}
}

func Int2Bytes(value int) []byte {
	return []byte(fmt.Sprintf("%d", value))
}
