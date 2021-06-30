package core

import (
	"github.com/dhindustries/graal/queue"
)

type runtime struct {
	r *queue.Invoker
}

func (rt *runtime) schedule(task func()) {
	rt.r.Q.Push(task)
}

func (rt *runtime) invoke(task func()) {
	rt.r.Call(task)
}
