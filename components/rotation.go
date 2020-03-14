package components

import (
	"sync"

	"github.com/dhindustries/graal"
	"github.com/go-gl/mathgl/mgl32"
)

type Rotated struct {
	rotation graal.Vec3
	matrix   graal.Mat4x4
	lock     sync.RWMutex
	valid    bool
}

func (component *Rotated) isValid() bool {
	component.lock.RLock()
	defer component.lock.RUnlock()
	return component.valid
}

func (component *Rotated) Rotation() graal.Vec3 {
	component.lock.RLock()
	defer component.lock.RUnlock()
	return component.rotation
}

func (component *Rotated) SetRotation(pos graal.Vec3) {
	component.lock.Lock()
	defer component.lock.Unlock()
	component.rotation = pos
	component.valid = false
}

func (component *Rotated) Transformation() graal.Mat4x4 {
	component.lock.Lock()
	defer component.lock.Unlock()
	if !component.valid {
		x := mgl32.HomogRotate3DX(component.rotation[0])
		y := mgl32.HomogRotate3DY(component.rotation[1])
		z := mgl32.HomogRotate3DZ(component.rotation[2])

		component.matrix = graal.Mat4x4(z.Mul4(y).Mul4(x))
		component.valid = true
	}
	return component.matrix
}
