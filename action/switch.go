package action

import (
	"sync/atomic"
	"time"
)

type Switch struct {
	Action Action
	Result bool
	value  atomic.Value
}

func (action *Switch) Run(t interface{}, dt time.Duration) bool {
	state := true
	if v := action.value.Load(); v != nil {
		state = v.(bool)
	}
	if state {
		return action.Action.Run(t, dt)
	}
	return action.Result
}

func (action *Switch) Enable() {
	action.value.Store(true)
}

func (action *Switch) Disable() {
	action.value.Store(false)
}
