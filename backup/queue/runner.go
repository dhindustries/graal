package queue

import "sync"

type Runner struct {
	Q *Queue
	b chan interface{}
	l sync.Mutex
}

func NewRunner(q *Queue) *Runner {
	return &Runner{Q: q}
}

func (runner *Runner) Start() {
	runner.l.Lock()
	if runner.b == nil {
		runner.b = make(chan interface{})
	}
	runner.l.Unlock()
	runner.Q.PullUntil(runner.b)
}

func (runner *Runner) Stop() {
	runner.l.Lock()
	defer runner.l.Unlock()
	runner.b <- nil
	close(runner.b)
	runner.b = nil
	runner.Q = nil
}
