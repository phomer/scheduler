package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/phomer/scheduler/accounts"
	"github.com/phomer/scheduler/comm"
	"github.com/phomer/scheduler/jobs"
	"github.com/phomer/scheduler/log"
	"github.com/phomer/scheduler/sig"
)

func main() {
	// Read in the config parms
	config := accounts.FindClientConfig()

	// Parse the args
	request := ParseArgs(config.Username)

	// Setup the SIGINT handler
	sig.Initialize()
	sig.Catch(syscall.SIGINT, StopStreaming)

	// Send the request
	response := comm.MakeRequest(config, request)

	// Stream the response back
	comm.DisplayStream(response)
}

type CliCommand int

var jobid CliCommand = 0
var cli_cmd string

var job_flags = map[string]bool{
	"tail":   true,
	"output": true,
	"status": true,
	"remove": true,
}

func ParseArgs(username string) *jobs.Request {

	// Add all of the flags as single arguments
	for value, _ := range job_flags {
		flag.Var(&jobid, value, "jobid")
	}

	flag.Parse()

	switch CountCommands() {
	case 0:
		return BuildExecuteRequest(username, os.Args)
	case 1:
		return BuildCmdRequest(username, os.Args)
	default:
		fmt.Println("Invalid arguments, too many flags")
	}

	return nil
}

func (cmd *CliCommand) String() string {
	return fmt.Sprintf("%d", cmd)
}

func (cmd *CliCommand) Set(value string) error {
	if *cmd != 0 {
		return errors.New("Duplicate Command")
	}

	integer, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*cmd = CliCommand(integer)

	return nil
}

func BuildCmdRequest(username string, args []string) *jobs.Request {
	return &jobs.Request{
		Username: username,
		Type:     cli_cmd,
		JobId:    int(jobid),
	}
}

func BuildExecuteRequest(username string, args []string) *jobs.Request {
	var start_scale, scale *jobs.TimeScale
	var start, cont int

	next := 1
	found := false

	found, start, start_scale = Frequency(args[next])
	if found {
		next++
		found, cont, scale = Frequency(args[next])
		if found {
			next++
		}
	}

	cmd := args[next]
	next++

	// Default the type to blank
	request_type := ""
	if next > 2 {
		request_type = "schedule"
	}

	return &jobs.Request{
		Username:   username,
		Type:       request_type,
		Cmd:        cmd,
		Args:       args[next:],
		Time:       time.Now().Unix(), // Local time for client
		Start:      start,
		StartScale: start_scale,
		Continue:   cont,
		Scale:      scale,
	}
}

func Frequency(value string) (bool, int, *jobs.TimeScale) {

	freqEx := regexp.MustCompile("([0-9]+)([^0-9]+)")

	list := freqEx.FindStringSubmatch(value)

	if len(list) == 3 {

		number, err := strconv.Atoi(list[1])
		if err != nil {
			// We can't validate this directly
			log.Fatal("Regex Failed", err)
		}

		scale := jobs.LookupTimeScale(list[2])
		if scale == nil {
			//Should assume its a command that starts with a number, not a time scale
			return false, 0, nil // Implicit defaults for the vars
		}

		return true, number, scale
	}
	return false, 0, nil // Implicit defaults for the vars
}

// Count up the number of command flags found in the args, for validation
func CountCommands() int {
	count := 0
	flag.Visit(func(value *flag.Flag) {

		// See if this is a flag
		_, ok := job_flags[value.Name]
		if ok {
			count++
		}

		// Side-effect: Save the last command
		cli_cmd = value.Name
	})
	return count
}

// Stop the client from streaming server job output
func StopStreaming() {
	fmt.Println("Shut her down, Clancy, she’s pumping mud")

	comm.StopStreaming = true

	go func() {
		// Give it a change to stop nicely, but if it's not happening ...
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}
