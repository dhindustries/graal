package components

import (
	"sync"

	"github.com/go-gl/mathgl/mgl64"
)

type scale struct {
	scale  mgl64.Vec3
	valid  bool
	matrix mgl64.Mat4
	lock   sync.RWMutex
}

func (component *scale) isValid() bool {
	component.lock.RLock()
	defer component.lock.RUnlock()
	return component.valid
}

func (component *scale) Scale() mgl64.Vec3 {
	component.lock.RLock()
	defer component.lock.RUnlock()
	return component.scale
}

func (component *scale) SetScale(pos mgl64.Vec3) {
	component.lock.Lock()
	defer component.lock.Unlock()
	component.scale = pos
	component.valid = false
}

func (component *scale) Transform() mgl64.Mat4 {
	component.lock.Lock()
	defer component.lock.Unlock()
	if !component.valid {
		component.matrix = mgl64.Scale3D(component.scale[0], component.scale[1], component.scale[2])
		component.valid = true
	}
	return component.matrix
}
