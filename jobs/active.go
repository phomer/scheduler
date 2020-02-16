package jobs

import (
	"os"
	"sync"
	"syscall"

	"github.com/phomer/scheduler/accounts"
	"github.com/phomer/scheduler/sig"
)

/* Threadsafe active children */
type ActiveJob struct {
	Cmd *Command

	IsRunning bool
	Pid       int
	Time      int64

	File   *os.File
	Offset int64
}

func NewActiveJob(account *accounts.Account, command *Command) *ActiveJob {

	sched := NewScheduled()
	command.JobId = sched.AllocateNewJobId(account.Username, command)

	return &ActiveJob{
		Cmd: command,
	}
}

// TODO: Build from an incoming command
func NewImmediateJob(account *accounts.Account, command *Command) *ActiveJob {
	//active := NewActive()
	sched := NewScheduled()
	command.JobId = sched.AllocateNewJobId(account.Username, command)

	return &ActiveJob{
		Cmd: command,
	}
}

type Active struct {
	Jobs map[int]*ActiveJob
	mux  sync.Mutex
}

var current *Active

func NewActive() *Active {
	if current != nil {
		return current
	}

	sig.Initialize()
	sig.Catch(syscall.SIGCHLD, UpdateJobStatus)

	return &Active{
		Jobs: make(map[int]*ActiveJob, 0),
	}
}

func (active *Active) IsActive(pid int) bool {

	active.mux.Lock()
	defer active.mux.Unlock()

	job, ok := active.Jobs[pid]
	if !ok {
		return false
	}

	return job.IsRunning
}

func (active *Active) AddJob(pid int, job *ActiveJob) {
	active.mux.Lock()

	active.Jobs[pid] = job

	active.mux.Unlock()
}

func (active *Active) FindJobStatus(username string, jobid int) *ActiveJob {
	active.mux.Lock()

	entry := active.find(username, jobid)

	active.mux.Unlock()

	return entry
}

func (active *Active) find(username string, jobid int) *ActiveJob {
	for _, entry := range active.Jobs {
		if entry.Cmd.Username == username && entry.Cmd.JobId == jobid {
			// TODO: Can't release this, it's not immutable
			return entry
		}
	}
	return nil
}

func UpdateJobStatus() {
	active := NewActive()

	active.mux.Lock()
	// Reap the child that died
	// Go through and delete entries in the overall list of children
	active.mux.Unlock()
}
