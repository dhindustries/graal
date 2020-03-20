package components

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type origin struct {
	origin mgl32.Vec3
	matrix mgl32.Mat4
	lock   sync.RWMutex
	valid  bool
}

func (component *origin) isValid() bool {
	component.lock.RLock()
	defer component.lock.RUnlock()
	return component.valid
}

func (component *origin) Origin() mgl32.Vec3 {
	component.lock.RLock()
	defer component.lock.RUnlock()
	return component.origin
}

func (component *origin) SetOrigin(pos mgl32.Vec3) {
	component.lock.Lock()
	defer component.lock.Unlock()
	component.origin = pos
	component.valid = false
}

func (component *origin) Transform() mgl32.Mat4 {
	component.lock.Lock()
	defer component.lock.Unlock()
	if !component.valid {
		component.matrix = mgl32.Mat4(mgl32.Translate3D(component.origin[0], component.origin[1], component.origin[2]))
		component.valid = true
	}
	return component.matrix
}
