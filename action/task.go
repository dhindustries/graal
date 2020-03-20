package action

import (
	"time"
)

type Action interface {
	Run(target interface{}, dt time.Duration) bool
}

type Event struct {
	Fn func()
}

func (action *Event) Run(t interface{}, dt time.Duration) bool {
	action.Fn()
	return true
}

type Delay struct {
	Duration time.Duration
	t        time.Duration
}

func (action *Delay) Run(t interface{}, dt time.Duration) bool {
	if action.t <= 0 {
		action.t = action.Duration
	}
	action.t -= dt
	return action.t <= 0
}

type Repeat struct {
	Action Action
	Times  int
	i      int
}

func (action *Repeat) Run(t interface{}, dt time.Duration) bool {
	if action.i <= 0 {
		action.i = action.Times
	}
	if action.Action.Run(t, dt) {
		action.i--
	}
	return action.i <= 0
}

type Timeout struct {
	Action   Action
	Duration time.Duration
	i        time.Duration
}

func (action *Timeout) Run(t interface{}, dt time.Duration) bool {
	if action.i <= 0 {
		action.i = action.Duration
	}
	action.i -= dt
	return action.i <= 0 || action.Run(t, dt)
}
