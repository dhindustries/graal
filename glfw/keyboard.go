package glfw

import (
	"sync"

	"github.com/dhindustries/graal"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type keyboard struct {
	cs, ps map[graal.Key]bool
	l      sync.RWMutex
}

func newKeyboard() *keyboard {
	return &keyboard{
		cs: make(map[graal.Key]bool),
		ps: make(map[graal.Key]bool),
	}
}

func (k *keyboard) isUp(api *graal.Api, key graal.Key) bool {
	k.l.RLock()
	defer k.l.RUnlock()
	return !k.cs[key]
}

func (k *keyboard) isDown(api *graal.Api, key graal.Key) bool {
	k.l.RLock()
	defer k.l.RUnlock()
	return k.cs[key]
}

func (k *keyboard) isPressed(api *graal.Api, key graal.Key) bool {
	k.l.RLock()
	defer k.l.RUnlock()

	return k.cs[key] && !k.ps[key]
}

func (k *keyboard) isReleased(api *graal.Api, key graal.Key) bool {
	k.l.RLock()
	defer k.l.RUnlock()

	return !k.cs[key] && k.ps[key]
}

func (k *keyboard) handle(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mod glfw.ModifierKey) {
	k.l.Lock()
	defer k.l.Unlock()
	k.cs[graal.Key(key)] = action != glfw.Release
}

func (k *keyboard) update(api *graal.Api) {
	k.l.Lock()
	defer k.l.Unlock()
	for key, state := range k.cs {
		k.ps[key] = state
	}
}
