package glfw

import (
	"sync"

	"github.com/dhindustries/graal"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type keyboard struct {
	currentState, previousState map[graal.Key]bool
	lock                        sync.RWMutex
}

func (keyboard *keyboard) IsUp(key graal.Key) bool {
	keyboard.lock.RLock()
	defer keyboard.lock.RUnlock()
	return !keyboard.currentState[key]
}

func (keyboard *keyboard) IsDown(key graal.Key) bool {
	keyboard.lock.RLock()
	defer keyboard.lock.RUnlock()
	return keyboard.currentState[key]
}

func (keyboard *keyboard) IsPressed(key graal.Key) bool {
	keyboard.lock.RLock()
	defer keyboard.lock.RUnlock()

	return keyboard.currentState[key] && !keyboard.previousState[key]
}

func (keyboard *keyboard) IsReleased(key graal.Key) bool {
	keyboard.lock.RLock()
	defer keyboard.lock.RUnlock()

	return !keyboard.currentState[key] && keyboard.previousState[key]
}

func (keyboard *keyboard) handle(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mod glfw.ModifierKey) {
	keyboard.lock.Lock()
	defer keyboard.lock.Unlock()
	keyboard.currentState[graal.Key(key)] = action != glfw.Release
}

func (keyboard *keyboard) update() {
	keyboard.lock.Lock()
	defer keyboard.lock.Unlock()
	for key, state := range keyboard.currentState {
		keyboard.previousState[key] = state
	}
}
