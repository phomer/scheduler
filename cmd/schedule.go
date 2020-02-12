package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/phomer/scheduler/accounts"
	"github.com/phomer/scheduler/comm"
	"github.com/phomer/scheduler/sig"
)

func main() {
	fmt.Println("schedule ", os.Args[1:])

	// Read in the config parms
	config := accounts.FindClientConfig()

	// Parse the args
	request := comm.NewRequest(os.Args)

	// Setup the SIGINT handler
	sig.Initialize()
	sig.Catch(syscall.SIGINT, StopStreaming)

	// Send the request
	response := comm.MakeRequest(config, request)

	// Stream the response back
	comm.DisplayStream(response)
}

func StopStreaming() {

	fmt.Println("Shut her down, Clancy, she’s pumping mud")

	// TODO: Should exit nicely, so that any incomming data is finsihed to the whole last string.
	os.Exit(0)
}
