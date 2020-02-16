package jobs

import (
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

	Next     int64
	Continue int
	Scale    TimeScale

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
		Next:     next,
	}
}

func NewScheduled() *Scheduled {
	user_map := &UserMap{
		Users: make(map[string]*CommandMap, 0),
	}

	return &Scheduled{
		Map: user_map,
		db:  datastore.NewDatabase("Jobs"),
	}
}

// Side-effect
func (cmap *CommandMap) update_jobid(command *Command) int {

	jobid := cmap.NextId
	cmap.NextId++
	command.JobId = jobid
	command.Filepath = OutputFilepath("data", command.Username, jobid)

	return jobid
}

func (cmap *CommandMap) add_command(command *Command) int {

	cmap.Commands[cmap.NextId] = command

	return cmap.update_jobid(command)
}

func (sched *Scheduled) find_user(username string) *CommandMap {
	user, ok := sched.Map.Users[username]
	if !ok {
		user = &CommandMap{
			Commands: make(map[int]*Command, 0),
			NextId:   101,
		}
		sched.Map.Users[username] = user
	}
	return user
}

func (sched *Scheduled) AllocateNewJobId(username string, cmd *Command) int {
	sched.mux.Lock()

	user := sched.find_user(username)
	jobid := user.update_jobid(cmd)

	sched.store()

	sched.mux.Unlock()

	return jobid
}

func (sched *Scheduled) AddScheduledCommand(username string, command *Command) int {
	sched.mux.Lock()

	user := sched.find_user(username)
	jobid := user.add_command(command)

	// Persist, no need for a lock, just one process has access
	sched.store()

	sched.mux.Unlock()

	return jobid
}

// Don't warn if the user or job isn't there
func (sched *Scheduled) RemoveUserCommand(username string, jobid int) int {
	sched.mux.Lock()

	user := sched.find_user(username)

	_, ok := user.Commands[jobid]
	if ok {
		user.Commands[jobid] = nil
	}

	// Persist, no need for a lock, just one process has access
	sched.store()

	sched.mux.Unlock()

	return jobid
}

// Reload the datastore
func (sched *Scheduled) Reload() {
	sched.mux.Lock()

	sched.load()

	sched.mux.Unlock()
}

// Go through every possible job and find the next one.
// TODO: Some form of indexing would be way better
func (sched *Scheduled) FindNext() (int64, []*Command) {
	sched.mux.Lock()

	sched.load()

	minimum := int64(0)
	set := make([]*Command, 0)

	// Full-table scan O(N) for all users, all jobs
	// We are okay, to do this in any order
	for _, list := range sched.Map.Users {

		for i := 0; i < len(list.Commands); i++ {

			start := list.Commands[i].Next
			if start > minimum {
				// Ignore
			} else {
				if start < minimum {
					// Reset the list
					set = make([]*Command, 0)
				}
				minimum = start
				set = append(set, list.Commands[i])
			}
		}
	}

	sched.mux.Unlock()

	return minimum, set
}

/* Wake up, and see what still needs to be done */
func ProcessSchedule(root interface{}) {

	sched := NewScheduled()
	//active := NewActive()

	for {
		now := time.Now().Unix() // Mark the time
		_ = now

		next, list := sched.FindNext()
		if next < 1 {
			// Order wasn't preserved in the find, so it doesn't matter here
			for _, command := range list {
				account := accounts.FindAccount(command.Username)
				job := NewActiveJob(account, command)
				Spawn(account, job)
			}
		} else {
			// TODO: Need to make sure whatever sleep method, can be interrupted.
			time.Sleep(time.Duration(next) * time.Second) // Can be interrupted by Signal
		}
	}
}

func (sched *Scheduled) load() {
	sched.Map = sched.db.Load(sched.Map).(*UserMap)
}

func (sched *Scheduled) store() {
	sched.db.Store(sched.Map)
}
