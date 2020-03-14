package graal

import "sync"

type Handle interface {
	Acquire()
	Release()
	Valid() bool
}

type BaseHandle struct {
	lock    sync.Mutex
	counter uint
	invalid bool
}

func (handle *BaseHandle) Acquire() {
	handle.lock.Lock()
	defer handle.lock.Unlock()
	if handle.invalid {
		panic("Acquire invalid handle")
	}
	handle.counter++
}

func (handle *BaseHandle) Release() {
	handle.lock.Lock()
	defer handle.lock.Unlock()
	if handle.invalid {
		panic("Release invalid handle")
	}
	if handle.counter > 0 {
		handle.counter--
	}
	if handle.counter == 0 {
		handle.invalid = true
	}
}

func (handle *BaseHandle) Valid() bool {
	handle.lock.Lock()
	defer handle.lock.Unlock()
	return !handle.invalid
}
