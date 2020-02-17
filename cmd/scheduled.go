package main

import (
	"fmt"

	"github.com/phomer/scheduler/accounts"
	"github.com/phomer/scheduler/comm"
)

func main() {
	daemonize()

	fmt.Println("Starting scheduled")
	fmt.Println("  Username:", accounts.Username())
	fmt.Println("  Hostname:", accounts.Hostname())

	// Start off the jobs queue

	// Start off the web service
	listen := comm.NewServer()
	listen.Start()
}

// Do all of the fun things necessary to daemonize a background service
func daemonize() {
	fmt.Println("Daemonize")

	// Fork a couple of times
	// Disconnect from the TTY
	// Reset stdin, stdout and stderr
}
