package action

import (
	"sync"
	"time"
)

type Sequence struct {
	d []Action
	i int
	l sync.RWMutex
}

func NewSequence() *Sequence {
	return &Sequence{d: make([]Action, 0)}
}

func (action *Sequence) Add(task Action) {
	action.l.Lock()
	defer action.l.Unlock()
	action.d = append(action.d, task)
}

func (action *Sequence) Run(t interface{}, dt time.Duration) bool {
	action.l.Lock()
	defer action.l.Unlock()
	if action.i >= len(action.d) {
		action.i = 0
	}
	if action.d[action.i].Run(t, dt) {
		action.i++
	}
	return action.i >= len(action.d)
}
