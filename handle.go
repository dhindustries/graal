package graal

import "sync"

type apiHandle interface {
	Acquire(handle interface{})
	Release(handle interface{})
}

type handler interface {
	Acquire()
	Release() bool
}

type disposer interface {
	Dispose()
}

type apiDisposer interface {
	Dispose(api Api)
}

func (api *apiAdapter) Acquire(handle interface{}) {
	Acquire(handle)
}

func (api *apiAdapter) Release(handle interface{}) {
	release(handle, api)
}

func Acquire(handle interface{}) {
	if h, ok := handle.(handler); ok {
		h.Acquire()
	}
}

func Release(handle interface{}) {
	release(handle, api)
}

func release(handle interface{}, api Api) {
	if h, ok := handle.(handler); ok {
		if h.Release() {
			dispose(handle, api)
		}
	} else {
		dispose(handle, api)
	}
}

func dispose(handle interface{}, api Api) {
	if d, ok := handle.(apiDisposer); ok {
		d.Dispose(api)
	} else if d, ok := handle.(disposer); ok {
		d.Dispose()
	}
}

type Handle struct {
	refs     uint32
	released bool
	lock     sync.Mutex
}

func (handle *Handle) Acquire() {
	handle.lock.Lock()
	defer handle.lock.Unlock()
	if handle.released {
		panic("cannot acquire released handle")
	}
	handle.refs++
}

func (handle *Handle) Release() bool {
	handle.lock.Lock()
	defer handle.lock.Unlock()
	if handle.released {
		panic("handle is already released")
	}
	if handle.refs > 0 {
		handle.refs--
	}
	if handle.refs == 0 {
		handle.released = true
	}
	return handle.released
}

type Bindable struct {
	Handle
	bound []handler
	lock  sync.Mutex
}

func (object *Bindable) Bind(handle interface{}) {
	if h, ok := handle.(handler); ok {
		object.lock.Lock()
		defer object.lock.Unlock()
		if object.bound == nil {
			object.bound = make([]handler, 1)
		}
		Acquire(h)
		object.bound = append(object.bound, h)
	}
}

func (object *Bindable) Unbind(handle interface{}) {
	object.lock.Lock()
	defer object.lock.Unlock()
	if h, ok := handle.(handler); object.bound != nil && ok {
		for i, bound := range object.bound {
			if bound == h {
				Release(h)
				object.bound = append(object.bound[:i], object.bound[i+1:]...)
				break
			}
		}
	}
}

func (object *Bindable) Acquire() {
	object.Handle.Acquire()
	object.lock.Lock()
	defer object.lock.Unlock()
	for _, bound := range object.bound {
		Acquire(bound)
	}
}

func (object *Bindable) Release() {
	object.lock.Lock()
	defer object.lock.Unlock()
	for _, bound := range object.bound {
		Release(bound)
	}
	object.Handle.Release()
}
