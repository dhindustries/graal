package action

import (
	"sync"
	"time"
)

type Queue struct {
	actions []Action
	l       sync.RWMutex
}

func (queue *Queue) Add(action Action) *Queue {
	queue.l.Lock()
	defer queue.l.Unlock()
	if queue.actions == nil {
		queue.actions = make([]Action, 0, 1)
	}
	queue.actions = append(queue.actions, action)
	return queue
}

func (queue *Queue) Run(t interface{}, dt time.Duration) bool {
	if action := queue.Top(); action != nil {
		if action.Run(t, dt) {
			queue.Pull()
		}
		return false
	}
	return true
}

func (queue *Queue) Empty() bool {
	queue.l.RLock()
	defer queue.l.RUnlock()
	return queue.actions == nil || len(queue.actions) == 0
}

func (queue *Queue) Top() Action {
	queue.l.RLock()
	defer queue.l.RUnlock()
	if queue.actions != nil && len(queue.actions) > 0 {
		return queue.actions[0]
	}
	return nil
}

func (queue *Queue) Pull() {
	queue.l.Lock()
	defer queue.l.Unlock()
	if queue.actions != nil && len(queue.actions) > 0 {
		queue.actions = queue.actions[1:]
	}
}

func (queue *Queue) Clear() {
	queue.l.Lock()
	defer queue.l.Unlock()
	queue.actions = nil
}
