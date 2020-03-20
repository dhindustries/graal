package components

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type rotation struct {
	rotation mgl32.Vec3
	matrix   mgl32.Mat4
	lock     sync.RWMutex
	valid    bool
}

func (component *rotation) isValid() bool {
	component.lock.RLock()
	defer component.lock.RUnlock()
	return component.valid
}

func (component *rotation) Rotation() mgl32.Vec3 {
	component.lock.RLock()
	defer component.lock.RUnlock()
	return component.rotation
}

func (component *rotation) SetRotation(pos mgl32.Vec3) {
	component.lock.Lock()
	defer component.lock.Unlock()
	component.rotation = pos
	component.valid = false
}

func (component *rotation) Transform() mgl32.Mat4 {
	component.lock.Lock()
	defer component.lock.Unlock()
	if !component.valid {
		x := mgl32.HomogRotate3DX(component.rotation[0])
		y := mgl32.HomogRotate3DY(component.rotation[1])
		z := mgl32.HomogRotate3DZ(component.rotation[2])

		component.matrix = mgl32.Mat4(z.Mul4(y).Mul4(x))
		component.valid = true
	}
	return component.matrix
}
