package components

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type Transformation struct {
	translation
	rotation
	scale
	origin
	matrix mgl32.Mat4
	lock   sync.Mutex
}

func (object *Transformation) isValid() bool {
	return object.translation.isValid() && object.rotation.isValid() && object.scale.isValid() && object.origin.isValid()
}

func (object *Transformation) Transform() mgl32.Mat4 {
	object.lock.Lock()
	defer object.lock.Unlock()
	if !object.isValid() {
		translation := mgl32.Mat4(object.translation.Transform())
		rotation := mgl32.Mat4(object.rotation.Transform())
		scale := mgl32.Mat4(object.scale.Transform())
		origin := mgl32.Mat4(object.origin.Transform())
		object.matrix = mgl32.Mat4(translation.Mul4(rotation).Mul4(scale).Mul4(origin))
	}
	return object.matrix
}
