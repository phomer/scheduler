package main

import (
	"fmt"
	"os"

	"github.com/phomer/scheduler/accounts"
	"github.com/phomer/scheduler/comm"
)

func main() {
	fmt.Println("schedule ", os.Args[1:])

	// Read in the config parms
	config := accounts.FindClientConfig()

	// Parse the args
	request := comm.NewRequest(os.Args)

	// Setup the SIGINT handler

	// Send the request
	comm.MakeRequest(config, request)

	// Stream the response back
}
