package jobs

import (
	"fmt"
	"sync"
	"time"

	"github.com/phomer/scheduler/accounts"
	"github.com/phomer/scheduler/datastore"
)

type Command struct {
	JobId    int
	Username string // Redundant, also in map key

	Cmd  string
	Args []string

	NextRun int64
	Pending bool

	Continue int
	Scale    *TimeScale

	Filepath string
}

type UserMap struct {
	Users map[string]*CommandMap
}

type CommandMap struct {
	Commands map[int]*Command
	NextId   int
}

type Scheduled struct {
	Map *UserMap
	db  *datastore.Database
	mux sync.Mutex
}

func NewCommand(request *Request) *Command {

	// User specified increment relative to them, convert that to be absolute for server
	curr := time.Now().Unix()
	base := curr + (curr - request.Time)

	next := int64(0)
	if request.StartScale != nil {
		next = AbsoluteUnixTime(base, request.Start, request.StartScale)
	}

	return &Command{
		Username: request.Username,
		Cmd:      request.Cmd,
		Args:     request.Args,
		NextRun:  next,
	}
}

var scheduler *Scheduled

func NewScheduled() *Scheduled {
	if scheduler != nil {
		return scheduler
	}

	user_map := &UserMap{
		Users: make(map[string]*CommandMap, 0),
	}

	scheduler = &Scheduled{
		Map: user_map,
		db:  datastore.NewDatabase("Jobs"),
	}

	scheduler.load()

	return scheduler
}

// Side-effect
func (cmap *CommandMap) update_jobid(command *Command) *Command {

	jobid := cmap.NextId

	command.JobId = jobid
	command.Filepath = OutputFilepath("data", command.Username, jobid)

	cmap.NextId++ // Side-effect

	return command
}

func (cmap *CommandMap) add_command(command *Command) *Command {

	command = cmap.update_jobid(command)
	command.Pending = true

	cmap.Commands[command.JobId] = command

	return command
}

func (sched *Scheduled) find_user(username string) *CommandMap {
	user, ok := sched.Map.Users[username]
	if !ok {
		user = &CommandMap{
			Commands: make(map[int]*Command, 0),
			NextId:   1,
		}
		sched.Map.Users[username] = user
	}
	return user
}

func (sched *Scheduled) AllocateNewJobId(username string, command *Command) *Command {
	sched.mux.Lock()
	defer sched.mux.Unlock()

	user := sched.find_user(username)
	command = user.update_jobid(command)

	sched.store()

	return command
}

func (sched *Scheduled) AddScheduledCommand(username string, command *Command) *Command {
	sched.mux.Lock()
	defer sched.mux.Unlock()

	user := sched.find_user(username)
	command = user.add_command(command)

	// Persist, no need for a lock, just one process has access
	sched.store()

	// Let the scheduler know that things have changed.
	Reschedule()

	return command
}

// Don't warn if the user or job isn't there
func (sched *Scheduled) RemoveUserCommand(username string, jobid int) int {
	sched.mux.Lock()
	defer sched.mux.Unlock()

	user := sched.find_user(username)

	_, ok := user.Commands[jobid]
	if ok {
		user.Commands[jobid] = nil
	}

	// Persist, no need for a lock, just one process has access
	sched.store()

	return jobid
}

// Reload the datastore
func (sched *Scheduled) Reload() {
	sched.mux.Lock()

	sched.load()

	sched.mux.Unlock()
}

func (sched *Scheduled) FindCommand(username string, jobid int) *Command {
	sched.mux.Lock()
	defer sched.mux.Unlock()

	user, ok := sched.Map.Users[username]
	if !ok {
		return nil
	}
	command, ok := user.Commands[jobid]
	if !ok {
		return nil
	}
	return command
}

func (sched *Scheduled) ResetCommand(command *Command) {

	if command.Scale != nil {
		fmt.Println("Found a rerun of ", command.Continue, command.Scale)
		base := time.Now().Unix()
		command.NextRun = AbsoluteUnixTime(base, command.Continue, command.Scale)
		command.Pending = true

	} else {
		command.NextRun = 0
		command.Pending = false
	}

	sched.mux.Lock()
	defer sched.mux.Unlock()

	fmt.Println("Updating the Schedule")
	sched.Map.Users[command.Username].Commands[command.JobId] = command

	sched.store()
}

// Go through every possible job and find the next one.
// TODO: Some form of indexing would be way better
func (sched *Scheduled) FindNext() (int64, []*Command) {

	sched.mux.Lock()
	defer sched.mux.Unlock()

	minimum := int64(0)
	set := make([]*Command, 0)
	found := false

	// Full-table scan O(N) for all users, all jobs
	// We are okay, to do this in any order
	for _, entry := range sched.Map.Users {
		for _, cmd := range entry.Commands {
			if !cmd.Pending {
				continue
			}

			// Seconds until we should start this, can be negative
			start := cmd.NextRun - time.Now().Unix()
			if !found {
				minimum = start
				set = append(set, cmd)
				found = true

			} else if start <= minimum {
				if start < minimum {
					// Reset the list
					set = make([]*Command, 0)
				}
				minimum = start
				set = append(set, cmd)
			}
		}
	}

	if !found {
		return int64(60), nil
	}

	return minimum, set
}

//
var scheduling = make(chan bool)

func Reschedule() {
	scheduling <- true
}

/* Wake up, and see what still needs to be done */
func ProcessSchedule() {
	sched := NewScheduled()

	for {
		fmt.Println("Scheduling work.")
		next, list := sched.FindNext()
		if next < 1 {
			fmt.Println("Executing jobs", list)

			// Order wasn't preserved in the find, so it doesn't matter here
			for _, command := range list {
				account := accounts.FindAccount(command.Username)

				job := NewActiveJob(command)

				// TODO: Mark the job as run, before we actually run it.
				sched.ResetCommand(command)

				Spawn(account, job)
			}
		} else {
			fmt.Println("Sleeping for", next, "until there is more work.")

			// Need to interupt sometimes.
			select {
			case <-scheduling:
				break
			case <-time.After(time.Duration(next) * time.Second):
				break
			}
		}
	}
}

func (sched *Scheduled) load() {
	sched.Map = sched.db.Load(sched.Map).(*UserMap)
}

func (sched *Scheduled) store() {
	sched.db.Store(sched.Map)
}
