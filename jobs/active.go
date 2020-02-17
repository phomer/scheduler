package jobs

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"syscall"

	"github.com/phomer/scheduler/accounts"
	"github.com/phomer/scheduler/log"
	"github.com/phomer/scheduler/sig"
)

/* Threadsafe active children */
type ActiveJob struct {
	Cmd *Command

	IsRunning bool
	Status    int
	Pid       int
	Start     int64

	File   *os.File
	Offset int64
}

func NewActiveJob(command *Command) *ActiveJob {
	return &ActiveJob{
		Cmd: command,
	}
}

func NewImmediateJob(account *accounts.Account, command *Command) *ActiveJob {

	sched := NewScheduled()
	command = sched.AllocateNewJobId(account.Username, command)

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

	current = &Active{
		Jobs: make(map[int]*ActiveJob, 0),
	}
	return current
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

	// Guard against coding failures.
	if pid == 0 {
		log.Fatal("Invalid Process", errors.New("Missing"), job)
	}

	active.mux.Lock()

	active.Jobs[pid] = job

	active.mux.Unlock()
}

func (active *Active) FindJobStatus(username string, jobid int) *ActiveJob {
	active.mux.Lock()

	entry := active.find(username, jobid)

	active.mux.Unlock()

	// Outside of the active lock
	if entry == nil {
		sched := NewScheduled()
		command := sched.FindCommand(username, jobid)
		if command != nil {
			entry = NewActiveJob(command)
		}
	}

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

func CheckStatus(pid int, job *ActiveJob) *ActiveJob {
	options := syscall.WNOHANG

	// Reap the child that died
	var status syscall.WaitStatus
	var usage syscall.Rusage

	_, err := syscall.Wait4(pid, &status, options, &usage)
	if err != nil {
		fmt.Println("Error for pid ", pid, err.(error).Error())
		job.IsRunning = false // Assume it is false for now.

	} else if status.Exited() {
		fmt.Println("Marking as Exited", pid)
		job.IsRunning = false
		job.Status = status.ExitStatus()
		job.Pid = pid

	} else {
		job.IsRunning = true
		job.Pid = pid
	}

	return job
}

func UpdateJobStatus() {
	fmt.Println("Waking up on a SIGCHLD")

	active := NewActive()

	active.mux.Lock()

	for pid, job := range active.Jobs {
		if job.IsRunning {
			active.Jobs[pid] = CheckStatus(pid, job)
		}
	}

	// Go through and delete entries in the overall list of children
	active.mux.Unlock()
}
