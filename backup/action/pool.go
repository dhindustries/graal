package action

import (
	"sync"
	"time"
)

type Pool struct {
	actions []Action
	l       sync.RWMutex
}

func (pool *Pool) Add(action Action) {
	pool.l.Lock()
	defer pool.l.Unlock()
	if pool.actions == nil {
		pool.actions = make([]Action, 0, 1)
	}
	pool.actions = append(pool.actions, action)
}

func (pool *Pool) Run(t interface{}, dt time.Duration) bool {
	pool.l.Lock()
	defer pool.l.Unlock()
	if pool.actions != nil {
		n := make([]Action, 0, len(pool.actions))
		for _, action := range pool.actions {
			if !action.Run(t, dt) {
				n = append(n, action)
			}
		}
		pool.actions = n
		return len(pool.actions) == 0
	}
	return true
}

type AsyncPool struct {
	Pool
}

func (pool *AsyncPool) Run(t interface{}, dt time.Duration) bool {
	pool.l.Lock()
	defer pool.l.Unlock()
	if pool.actions != nil {
		var wg sync.WaitGroup
		data := make(chan Action, len(pool.actions))

		for _, action := range pool.actions {
			wg.Add(1)
			go func(action Action) {
				if !action.Run(t, dt) {
					data <- action
				}
				wg.Done()
			}(action)
		}
		go func() {
			wg.Wait()
			close(data)
		}()

		pool.actions = []Action{}
		for action := range data {
			pool.actions = append(pool.actions, action)
		}

		return len(pool.actions) == 0
	}
	return true
}
