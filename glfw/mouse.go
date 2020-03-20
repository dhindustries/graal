package glfw

import (
	"sync"

	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/dhindustries/graal"
	"github.com/go-gl/mathgl/mgl32"
)

type mouse struct {
	ps, cs map[graal.MouseButton]bool
	x, y   float32
	l      sync.RWMutex
}

func newMouse() *mouse {
	return &mouse{
		ps: make(map[graal.MouseButton]bool),
		cs: make(map[graal.MouseButton]bool),
	}
}

func (m *mouse) isUp(api *graal.Api, b graal.MouseButton) bool {
	m.l.RLock()
	defer m.l.RUnlock()
	return !m.cs[b]
}

func (m *mouse) isDown(api *graal.Api, b graal.MouseButton) bool {
	m.l.RLock()
	defer m.l.RUnlock()
	return m.cs[b]
}

func (m *mouse) isPressed(api *graal.Api, b graal.MouseButton) bool {
	m.l.RLock()
	defer m.l.RUnlock()

	return m.cs[b] && !m.ps[b]
}

func (m *mouse) isReleased(api *graal.Api, b graal.MouseButton) bool {
	m.l.RLock()
	defer m.l.RUnlock()

	return !m.cs[b] && m.ps[b]
}

func (m *mouse) cursor(api *graal.Api) mgl32.Vec2 {
	m.l.RLock()
	defer m.l.RUnlock()

	return mgl32.Vec2{m.x, m.y}
}

func (m *mouse) handleCursor(w *glfw.Window, xpos float64, ypos float64) {
	m.l.Lock()
	defer m.l.Unlock()
	ww, wh := w.GetSize()
	m.x = float32(xpos) / float32(ww)
	m.y = float32(ypos) / float32(wh)
}

func (m *mouse) handleButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	m.l.Lock()
	defer m.l.Unlock()
	m.cs[graal.MouseButton(button)] = action != glfw.Release
}

func (m *mouse) update(api *graal.Api) {
	m.l.Lock()
	defer m.l.Unlock()
	for key, state := range m.cs {
		m.ps[key] = state
	}
}
