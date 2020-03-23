package components

import (
	"sync"

	"github.com/go-gl/mathgl/mgl64"
)

type Transformation struct {
	translation
	rotation
	scale
	origin
	matrix mgl64.Mat4
	lock   sync.Mutex
}

func (object *Transformation) isValid() bool {
	return object.translation.isValid() && object.rotation.isValid() && object.scale.isValid() && object.origin.isValid()
}

func (object *Transformation) Transform() mgl64.Mat4 {
	object.lock.Lock()
	defer object.lock.Unlock()
	if !object.isValid() {
		translation := object.translation.Transform()
		rotation := object.rotation.Transform()
		scale := object.scale.Transform()
		origin := object.origin.Transform()
		object.matrix = translation.Mul4(rotation).Mul4(scale).Mul4(origin)
	}
	return object.matrix
}
