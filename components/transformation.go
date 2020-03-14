package components

import (
	"sync"

	"github.com/dhindustries/graal"
	"github.com/go-gl/mathgl/mgl32"
)

type Transformed struct {
	Translated
	Rotated
	Scaled
	matrix graal.Mat4x4
	lock   sync.Mutex
}

func (object *Transformed) Transformation() graal.Mat4x4 {
	object.lock.Lock()
	defer object.lock.Unlock()
	if !(object.Translated.isValid() && object.Rotated.isValid() && object.Scaled.isValid()) {
		translation := mgl32.Mat4(object.Translated.Transformation())
		rotation := mgl32.Mat4(object.Rotated.Transformation())
		scale := mgl32.Mat4(object.Scaled.Transformation())
		object.matrix = graal.Mat4x4(translation.Mul4(rotation).Mul4(scale))
	}
	return object.matrix
}
