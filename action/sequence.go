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

func (action *Sequence) Add(task Action) {
	action.l.Lock()
	defer action.l.Unlock()
	action.d = append(action.d, task)
}

func (action *Sequence) Run(t interface{}, dt time.Duration) bool {
	action.l.Lock()
	defer action.l.Unlock()
	if action.d == nil || len(action.d) == 0 {
		return true
	}
	if action.i >= len(action.d) {
		action.i = 0
	}
	if action.d[action.i].Run(t, dt) {
		action.i++
	}
	return action.i >= len(action.d)
}
