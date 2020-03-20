package queue

import "sync"

type TaskFn = func()

type Queue struct {
	l sync.Mutex
	t chan TaskFn
}

func (queue *Queue) init() {
	queue.l.Lock()
	defer queue.l.Unlock()
	if queue.t == nil {
		queue.t = make(chan TaskFn, 4096)
	}
}

func (queue *Queue) Push(task TaskFn) {
	if queue != nil {
		queue.init()
		queue.t <- task
	} else {
		task()
	}
}

func (queue *Queue) Pull() {
	if queue != nil {
		queue.init()
		t := <-queue.t
		t()
	}
}

func (queue *Queue) PullAll() {
	if queue != nil {
		queue.init()
	loop:
		for {
			select {
			case t := <-queue.t:
				t()
			default:
				break loop
			}
		}
	}
}

func (queue *Queue) PullUntil(brk chan interface{}) {
	if queue != nil {
		queue.init()
	loop:
		for {
			select {
			case t := <-queue.t:
				t()
			case <-brk:
				break loop
			}
		}
	}
}
