package jobs

import (
	"sync"

	"github.com/phomer/scheduler/datastore"
)

type Command struct {
	Cmd  string
	Args []string
}

type CommandMap struct {
	Commands map[string]*Command
}

type Scheduled struct {
	Map *CommandMap
	db  *datastore.Database
	mux sync.Mutex
}

func NewCommand(cmd string, args []string) *Command {
	return &Command{
		Cmd:  cmd,
		Args: args,
	}
}

func NewScheduled() *Scheduled {
	command_map := &CommandMap{
		Commands: make(map[string]*Command, 0),
	}

	return &Scheduled{
		Map: command_map,
		db:  datastore.NewDatabase("Jobs"),
	}
}

func (sched *Scheduled) AddCommand(username string, cmd string, args []string) {
	sched.mux.Lock()

	command := NewCommand(cmd, args)
	sched.Map.Commands[username] = command

	sched.mux.Unlock()
}

/* Wake up, and see what still needs to be done */
