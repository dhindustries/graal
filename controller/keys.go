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
	c chan interface{}
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

func (k *Keys) Track(freq float32) {
	k.l.Lock()
	if k.c != nil {
		k.l.Unlock()
		return
	}
	k.c = make(chan interface{})
	k.l.Unlock()
	t := time.NewTicker(time.Second / time.Duration(freq))
	go func() {
		p := time.Now()
	loop:
		for {
			select {
			case <-k.c:
				t.Stop()
				break loop
			case t := <-t.C:
				k.Update(t.Sub(p))
				p = t
			}
		}
		k.l.Lock()
		defer k.l.Unlock()
		close(k.c)
		t.Stop()
		k.c = nil
	}()
}

func (k *Keys) Stop() {
	k.l.Lock()
	defer k.l.Unlock()
	if k.c != nil {
		k.c <- true
	}
}

func (k *Keys) Dispose() {
	k.l.Lock()
	defer k.l.Unlock()
	if k.c != nil {
		k.c <- true
	}
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
