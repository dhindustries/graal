package components

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type translation struct {
	position mgl32.Vec3
	matrix   mgl32.Mat4
	lock     sync.RWMutex
	valid    bool
}

func (object *translation) isValid() bool {
	object.lock.RLock()
	defer object.lock.RUnlock()
	return object.valid
}

func (object *translation) Position() mgl32.Vec3 {
	object.lock.RLock()
	defer object.lock.RUnlock()
	return object.position
}

func (object *translation) SetPosition(pos mgl32.Vec3) {
	object.lock.Lock()
	defer object.lock.Unlock()
	object.position = pos
	object.valid = false
}

func (object *translation) Transform() mgl32.Mat4 {
	object.lock.Lock()
	defer object.lock.Unlock()
	if !object.valid {
		object.matrix = mgl32.Mat4(mgl32.Translate3D(object.position[0], object.position[1], object.position[2]))
		object.valid = true
	}
	return object.matrix
}
