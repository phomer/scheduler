package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
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
	request := ParseArgs()

	// Setup the SIGINT handler
	sig.Initialize()
	sig.Catch(syscall.SIGINT, StopStreaming)

	// Send the request
	response := comm.MakeRequest(config, request)

	// Stream the response back
	comm.DisplayStream(response)
}

type CliCommand int

var jobid CliCommand
var cli_cmd string

func ParseArgs() *comm.Request {
	flag.Var(&jobid, "tail", "jobid for scheduled process")
	flag.Var(&jobid, "output", "jobid for scheduled process")
	flag.Var(&jobid, "status", "jobid for scheduled process")
	flag.Var(&jobid, "remove", "jobid for scheduled process")

	flag.Parse()

	switch CountCommands() {
	case 0:
		return BuildExecuteRequest(os.Args)
	case 1:
		return BuildCmdRequest(os.Args)
	default:
		fmt.Println("Invalid arguments, too many flags")
	}

	return nil
}

func (cmd *CliCommand) String() string {
	return fmt.Sprintf("%d", cmd)
}

func (cmd *CliCommand) Set(value string) error {
	fmt.Println("Found a flag " + value)

	if *cmd != 0 {
		return errors.New("Duplicate Command")
	}
	integer, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	cmd = (*CliCommand)(&integer)

	return nil
}

func BuildCmdRequest(args []string) *comm.Request {
	return &comm.Request{
		Type:  cli_cmd,
		JobId: int(jobid),
	}
}

func BuildExecuteRequest(args []string) *comm.Request {
	var scale, cont_scale string
	var start, cont int

	next := 1
	found := false

	found, start, scale = Frequency(args[next])
	if found {
		next++
		found, cont, cont_scale = Frequency(args[next])
		if found {
			next++
		}
	}

	cmd := args[next]
	next++

	return &comm.Request{
		Command:       cmd,
		Args:          args[next:],
		Start:         start,
		StartScale:    scale,
		Continue:      cont,
		ContinueScale: cont_scale,
	}
}

func Frequency(value string) (bool, int, string) {

	freqEx := regexp.MustCompile("([0-9]+)([^0-9]+)")

	list := freqEx.FindStringSubmatch(value)

	if len(list) == 2 {
		number, err := strconv.Atoi(list[0])
		if err != nil {
			log.Fatal("Regex Failed", err)
		}

		return true, number, list[1]
	}
	return false, 0, "" // Implicit defaults for the vars
}

func CountCommands() int {
	count := 0
	flag.Visit(func(value *flag.Flag) {
		fmt.Println("Visiting " + value.Name)

		// TODO: Shouldn't be this redundant.
		if value.Name == "tail" || value.Name == "output" ||
			value.Name == "status" || value.Name == "remove" {

			count++
		}

		// Side-effect: Save the last command
		cli_cmd = value.Name
	})
	return count
}

func StopStreaming() {

	fmt.Println("Shut her down, Clancy, she’s pumping mud")

	// TODO: Should exit nicely, so that any incomming data is finsihed to the whole last string.
	os.Exit(0)
}
