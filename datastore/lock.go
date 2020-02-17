package datastore

import (
	"github.com/juju/fslock"
)

type FileLock struct {
	lock *fslock.Lock
}

func Lock(filename string) *FileLock {
	return &FileLock{
		lock: fslock.New(filename),
	}
}

func (file_lock *FileLock) Unlock() {
	file_lock.lock.Unlock()
}
