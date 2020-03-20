package components

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type scale struct {
	scale  mgl32.Vec3
	valid  bool
	matrix mgl32.Mat4
	lock   sync.RWMutex
}

func (component *scale) isValid() bool {
	component.lock.RLock()
	defer component.lock.RUnlock()
	return component.valid
}

func (component *scale) Scale() mgl32.Vec3 {
	component.lock.RLock()
	defer component.lock.RUnlock()
	return component.scale
}

func (component *scale) SetScale(pos mgl32.Vec3) {
	component.lock.Lock()
	defer component.lock.Unlock()
	component.scale = pos
	component.valid = false
}

func (component *scale) Transform() mgl32.Mat4 {
	component.lock.Lock()
	defer component.lock.Unlock()
	if !component.valid {
		component.matrix = mgl32.Mat4(mgl32.Scale3D(component.scale[0], component.scale[1], component.scale[2]))
		component.valid = true
	}
	return component.matrix
}
