package graal

import (
	"sync"
)

type Queue struct {
	tasks chan func()
	done  chan bool
	lock  sync.Mutex
}

func (queue *Queue) init() {
	queue.lock.Lock()
	defer queue.lock.Unlock()
	if queue.tasks == nil {
		queue.tasks = make(chan func(), 1)
		queue.done = make(chan bool, 1)
	}
}

func (queue *Queue) Pull() {
	if queue != nil {
		queue.init()
	loop:
		for true {
			select {
			case task := <-queue.tasks:
				task()
			default:
				break loop
			}
		}
	}
}

func (queue *Queue) Run() {
	if queue != nil {
		queue.init()
	loop:
		for true {
			select {
			case task := <-queue.tasks:
				task()
			case <-queue.done:
				break loop
			}
		}
	}
}

func (queue *Queue) Break() {
	if queue != nil {
		queue.done <- true
	}
}

func (queue *Queue) Push(task func()) {
	if queue != nil {
		queue.init()
		queue.tasks <- task
	} else {
		task()
	}
}

func (queue *Queue) Exec(task func() error) error {
	result := make(chan error, 1)
	queue.Push(func() {
		result <- task()
	})
	return <-result
}
