package component

import (
	"sync"
	"sync/atomic"

	"github.com/go-gl/mathgl/mgl64"
)

type Transformed struct {
	position       atomic.Value
	rotation       atomic.Value
	scale          atomic.Value
	origin         atomic.Value
	positionValid  bool
	rotationValid  bool
	scaleValid     bool
	originValid    bool
	positionMatrix mgl64.Mat4
	rotationMatrix mgl64.Mat4
	scaleMatrix    mgl64.Mat4
	originMatrix   mgl64.Mat4
	matrix         mgl64.Mat4
	mutex          sync.RWMutex
}

func (comp *Transformed) Position() mgl64.Vec3 {
	if v, ok := comp.position.Load().(mgl64.Vec3); ok {
		return v
	}
	return mgl64.Vec3{0, 0, 0}
}

func (comp *Transformed) SetPosition(value mgl64.Vec3) {
	comp.mutex.RLock()
	defer comp.mutex.RUnlock()
	comp.position.Store(value)
	comp.positionValid = false
}

func (comp *Transformed) Rotation() mgl64.Vec3 {
	if v, ok := comp.rotation.Load().(mgl64.Vec3); ok {
		return v
	}
	return mgl64.Vec3{0, 0, 0}
}

func (comp *Transformed) SetRotation(value mgl64.Vec3) {
	comp.mutex.RLock()
	defer comp.mutex.RUnlock()
	comp.rotation.Store(value)
	comp.rotationValid = false
}

func (comp *Transformed) Scale() mgl64.Vec3 {
	if v, ok := comp.scale.Load().(mgl64.Vec3); ok {
		return v
	}
	return mgl64.Vec3{1, 1, 1}
}

func (comp *Transformed) SetScale(value mgl64.Vec3) {
	comp.mutex.RLock()
	defer comp.mutex.RUnlock()
	comp.scale.Store(value)
	comp.scaleValid = false
}

func (comp *Transformed) Origin() mgl64.Vec3 {
	if v, ok := comp.origin.Load().(mgl64.Vec3); ok {
		return v
	}
	return mgl64.Vec3{0, 0, 0}
}

func (comp *Transformed) SetOrigin(value mgl64.Vec3) {
	comp.mutex.RLock()
	defer comp.mutex.RUnlock()
	comp.origin.Store(value)
	comp.originValid = false
}

func (comp *Transformed) Transformation() mgl64.Mat4 {
	comp.mutex.Lock()
	defer comp.mutex.Unlock()
	if !comp.positionValid {
		v := comp.Position()
		comp.positionMatrix = mgl64.Translate3D(v[0], v[1], v[2])
	}
	if !comp.rotationValid {
		v := comp.Rotation()
		x := mgl64.HomogRotate3DX(v[0])
		y := mgl64.HomogRotate3DY(v[1])
		z := mgl64.HomogRotate3DZ(v[2])
		comp.rotationMatrix = z.Mul4(y).Mul4(x)
	}
	if !comp.scaleValid {
		v := comp.Scale()
		comp.scaleMatrix = mgl64.Scale3D(v[0], v[1], v[2])
	}
	if !comp.originValid {
		v := comp.Origin()
		comp.originMatrix = mgl64.Translate3D(v[0], v[1], v[2])
	}
	if !(comp.positionValid && comp.rotationValid && comp.scaleValid && comp.originValid) {
		comp.matrix = comp.positionMatrix.Mul4(comp.rotationMatrix).Mul4(comp.scaleMatrix).Mul4(comp.originMatrix)
	}
	comp.positionValid = true
	comp.rotationValid = true
	comp.scaleValid = true
	comp.originValid = true
	return comp.matrix
}
