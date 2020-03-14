package components

import (
	"sync"

	"github.com/dhindustries/graal"
	"github.com/go-gl/mathgl/mgl32"
)

type Translated struct {
	position graal.Vec3
	matrix   graal.Mat4x4
	lock     sync.RWMutex
	valid    bool
}

func (object *Translated) isValid() bool {
	object.lock.RLock()
	defer object.lock.RUnlock()
	return object.valid
}

func (object *Translated) Position() graal.Vec3 {
	object.lock.RLock()
	defer object.lock.RUnlock()
	return object.position
}

func (object *Translated) SetPosition(pos graal.Vec3) {
	object.lock.Lock()
	defer object.lock.Unlock()
	object.position = pos
	object.valid = false
}

func (object *Translated) Transformation() graal.Mat4x4 {
	object.lock.Lock()
	defer object.lock.Unlock()
	if !object.valid {
		object.matrix = graal.Mat4x4(mgl32.Translate3D(object.position[0], object.position[1], object.position[2]))
		object.valid = true
	}
	return object.matrix
}
