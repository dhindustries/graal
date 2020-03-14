package components

import (
	"sync"

	"github.com/dhindustries/graal"
	"github.com/go-gl/mathgl/mgl32"
)

type Scaled struct {
	scale  graal.Vec3
	valid  bool
	matrix graal.Mat4x4
	lock   sync.RWMutex
}

func (component *Scaled) isValid() bool {
	component.lock.RLock()
	defer component.lock.RUnlock()
	return component.valid
}

func (component *Scaled) Scale() graal.Vec3 {
	component.lock.RLock()
	defer component.lock.RUnlock()
	return component.scale
}

func (component *Scaled) SetScale(pos graal.Vec3) {
	component.lock.Lock()
	defer component.lock.Unlock()
	component.scale = pos
	component.valid = false
}

func (component *Scaled) Transformation() graal.Mat4x4 {
	component.lock.Lock()
	defer component.lock.Unlock()
	if !component.valid {
		component.matrix = graal.Mat4x4(mgl32.Scale3D(component.scale[0], component.scale[1], component.scale[2]))
		component.valid = true
	}
	return component.matrix
}
