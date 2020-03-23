package components

import (
	"sync"

	"github.com/go-gl/mathgl/mgl64"
)

type translation struct {
	position mgl64.Vec3
	matrix   mgl64.Mat4
	lock     sync.RWMutex
	valid    bool
}

func (object *translation) isValid() bool {
	object.lock.RLock()
	defer object.lock.RUnlock()
	return object.valid
}

func (object *translation) Position() mgl64.Vec3 {
	object.lock.RLock()
	defer object.lock.RUnlock()
	return object.position
}

func (object *translation) SetPosition(pos mgl64.Vec3) {
	object.lock.Lock()
	defer object.lock.Unlock()
	object.position = pos
	object.valid = false
}

func (object *translation) Transform() mgl64.Mat4 {
	object.lock.Lock()
	defer object.lock.Unlock()
	if !object.valid {
		object.matrix = mgl64.Translate3D(object.position[0], object.position[1], object.position[2])
		object.valid = true
	}
	return object.matrix
}
