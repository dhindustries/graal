package graal

import "sync"

type Queue struct {
	tasks chan func() error
	lock  sync.Mutex
}

func (queue *Queue) init() {
	queue.lock.Lock()
	defer queue.lock.Unlock()
	if queue.tasks == nil {
		queue.tasks = make(chan func() error, 128)
	}
}

func (queue *Queue) invoke(task func() error) {
	task()
}

func (queue *Queue) Pull() {
	queue.init()
	done := false
	for !done {
		select {
		case task := <-queue.tasks:
			queue.invoke(task)
		default:
			done = true
		}
	}
}

func (queue *Queue) Push(task func() error) {
	queue.init()
	queue.tasks <- task
}
