package controller

import (
	"sync"
	"time"

	"github.com/dhindustries/graal"
)

type Keys struct {
	b graal.Keyboard
	m map[graal.Key]func()
	l sync.RWMutex
}

func NewKeys(keyboard graal.Keyboard) *Keys {
	return &Keys{b: keyboard, m: make(map[graal.Key]func())}
}

func (k *Keys) Bind(key graal.Key, fn func()) {
	k.l.Lock()
	defer k.l.Unlock()
	k.m[key] = fn
}

func (k *Keys) Unbind(key graal.Key) {
	k.l.Lock()
	defer k.l.Unlock()
	delete(k.m, key)
}

func (k *Keys) Update(dt time.Duration) {
	k.l.Lock()
	defer k.l.Unlock()
	for key, fn := range k.m {
		if k.b.IsDown(key) {
			go fn()
		}
	}
}
