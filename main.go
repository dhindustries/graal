package graal

import (
	"runtime"
	"sync"
)

var mainQueue Queue

type apiMain interface {
	Invoke(fn func(api Api))
	TryInvoke(fn func(api Api) error) error
}

func LockMainThread() {
	mainQueue.Start()
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		mainQueue.Pull()
	}()
}

func UnlockMainThread() {
	mainQueue.Stop()
}

func (api *apiAdapter) Invoke(fn func(api Api)) {
	if api.immediate {
		fn(api)
	} else {
		fork := &apiAdapter{api.proto, true}
		var wg sync.WaitGroup
		wg.Add(1)
		mainQueue.Enqueue(func() {
			fn(fork)
			wg.Done()
		})
		wg.Wait()
	}
}

func (api *apiAdapter) TryInvoke(fn func(api Api) error) error {
	if api.immediate {
		return fn(api)
	} else {
		fork := &apiAdapter{api.proto, true}
		res := make(chan error)
		defer close(res)
		mainQueue.Enqueue(func() {
			res <- fn(fork)
		})
		return <-res
	}
}
