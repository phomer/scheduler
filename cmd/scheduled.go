package main

import (
	"fmt"

	"github.com/phomer/scheduler/comm"
)

func main() {
	fmt.Println("Starting the scheduled Daemon")

	daemonize()

	// Setup handlers for SIGINT, and SIGHUP
	// Start off the jobs queue

	// Start off the web service
	listen := comm.NewServer()
	listen.Start()
}

// Do all of the fun things necessary to daemonize a background service
func daemonize() {
	fmt.Println("Fixing the daemon")
	// Fork a couple of times
	// Disconnect from the TTY
	// Reset stdin, stdout and stderr
}
