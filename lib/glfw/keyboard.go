package glfw

import (
	"sync"

	"github.com/dhindustries/graal"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type keyboard struct {
	window           *glfw.Window
	state            map[graal.Key]bool
	prevState        map[graal.Key]bool
	stateMutex       sync.RWMutex
	charChannels     []chan rune
	charMutex        sync.RWMutex
	prevCharCallback glfw.CharCallback
	prevKeyCallback  glfw.KeyCallback
}

func (keyboard *keyboard) init(window *glfw.Window) {
	keyboard.window = window
	keyboard.state = make(map[graal.Key]bool, int(glfw.KeyLast))
	keyboard.prevState = make(map[graal.Key]bool, int(glfw.KeyLast))
	keyboard.charChannels = make([]chan rune, 0)
	keyboard.prevCharCallback = window.SetCharCallback(keyboard.handleChar)
	keyboard.prevKeyCallback = window.SetKeyCallback(keyboard.handleKey)
}

func (keyboard *keyboard) finish() {
	keyboard.stateMutex.Lock()
	defer keyboard.stateMutex.Unlock()
	keyboard.charMutex.Lock()
	defer keyboard.charMutex.Unlock()
	keyboard.window.SetCharCallback(keyboard.prevCharCallback)
	keyboard.window = nil
	keyboard.prevCharCallback = nil
	keyboard.prevKeyCallback = nil
	keyboard.prevState = nil
	keyboard.state = nil
	for _, ch := range keyboard.charChannels {
		close(ch)
	}
	keyboard.charChannels = nil
}

func (keyboard *keyboard) update() {
	keyboard.stateMutex.Lock()
	defer keyboard.stateMutex.Unlock()
	for key := graal.KeySpace; key <= graal.KeyLast; key++ {
		keyboard.prevState[key] = keyboard.state[key]
		keyboard.state[key] = keyboard.window.GetKey(glfw.Key(key)) != glfw.Release
	}
}

func (keyboard *keyboard) handleChar(window *glfw.Window, char rune) {
	keyboard.sendChar(char)
	if keyboard.prevCharCallback != nil {
		keyboard.prevCharCallback(window, char)
	}
}

func (keyboard *keyboard) handleKey(
	window *glfw.Window,
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mods glfw.ModifierKey,
) {
	if action == glfw.Press || action == glfw.Repeat {
		switch key {
		case glfw.KeyEnter:
			keyboard.sendChar('\n')
		case glfw.KeyBackspace:
			keyboard.sendChar('\x08')
		case glfw.KeyDelete:
			keyboard.sendChar('\x7F')
		}
	}
	if keyboard.prevKeyCallback != nil {
		keyboard.prevKeyCallback(window, key, scancode, action, mods)
	}
}

func (keyboard *keyboard) isKeyDown(api graal.Api, key graal.Key) bool {
	keyboard.stateMutex.RLock()
	defer keyboard.stateMutex.RUnlock()
	return keyboard.state[key]
}

func (keyboard *keyboard) isKeyUp(api graal.Api, key graal.Key) bool {
	keyboard.stateMutex.RLock()
	defer keyboard.stateMutex.RUnlock()
	return !keyboard.state[key]
}

func (keyboard *keyboard) isKeyPressed(api graal.Api, key graal.Key) bool {
	keyboard.stateMutex.RLock()
	defer keyboard.stateMutex.RUnlock()
	return keyboard.state[key] && !keyboard.prevState[key]
}

func (keyboard *keyboard) isKeyReleased(api graal.Api, key graal.Key) bool {
	keyboard.stateMutex.RLock()
	defer keyboard.stateMutex.RUnlock()
	return !keyboard.state[key] && keyboard.prevState[key]
}

func (keyboard *keyboard) input(api graal.Api) (<-chan rune, func()) {
	ch := keyboard.addCharChannel()
	closer := func() {
		keyboard.removeCharChannel(ch)
	}
	return ch, closer
}

func (keyboard *keyboard) addCharChannel() chan rune {
	keyboard.charMutex.Lock()
	defer keyboard.charMutex.Unlock()
	ch := make(chan rune, 256)
	keyboard.charChannels = append(keyboard.charChannels, ch)
	return ch
}

func (keyboard *keyboard) removeCharChannel(ch chan rune) {
	keyboard.charMutex.Lock()
	defer keyboard.charMutex.Unlock()
	for i, c := range keyboard.charChannels {
		if c == ch {
			keyboard.charChannels = append(
				keyboard.charChannels[:i],
				keyboard.charChannels[i+1:]...,
			)
			close(ch)
			return
		}
	}
}

func (keyboard *keyboard) sendChar(char rune) {
	keyboard.charMutex.RLock()
	defer keyboard.charMutex.RUnlock()
	for _, ch := range keyboard.charChannels {
		go func() {
			ch <- char
		}()
	}
}
