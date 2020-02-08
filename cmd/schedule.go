package main

import (
	"fmt"
	"os"

	"github.com/phomer/scheduler/comm"
)

func main() {
	fmt.Println("schedule ", os.Args[1:])

	// Parse the args
	// Send the request
	// Setup the SIGINT handler
	// Stream the response back
	comm.Request("Hello There?")
}
