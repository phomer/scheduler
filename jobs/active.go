package jobs

import "sync"

/* Threadsafe active children */
type JobStatus struct {
	Cmd    *Command
	Status bool
	Pid    uint
}

type Active struct {
	Child map[uint]*JobStatus
	mux   sync.Mutex
}

func NewActive() *Active {
	return &Active{
		Child: make(map[uint]*JobStatus, 0),
	}
}

func (active *Active) AddJob(pid uint) {
	active.mux.Lock()

	active.Child[pid] = NewJobStatus()

	active.mux.Unlock()
}

func NewJobStatus() *JobStatus {
	return &JobStatus{}
}
