package core

import (
	"sync"

	"github.com/dhindustries/graal"
)

type handle struct {
	counter uint
	lock    sync.Mutex
	invalid bool
}

func (h *handle) Acquire() {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.invalid {
		panic("Acquire invalid handle")
	}
	h.counter++
}

func (h *handle) Release() {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.invalid {
		panic("Release invalid handle")
	}
	if h.counter > 0 {
		h.counter--
	}
	if h.counter == 0 {
		h.invalid = true
	}
}

func (h *handle) IsValid() bool {
	h.lock.Lock()
	defer h.lock.Unlock()
	return !h.invalid
}

func newHandle(*graal.Api) graal.Handle {
	return &handle{counter: 1}
}
