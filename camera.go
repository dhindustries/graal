package graal

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera interface {
	View() Mat4x4
	Projection() Mat4x4
}

type BaseCamera struct {
	position, lookAt, up     Vec3
	far, near                float32
	view, proj               Mat4x4
	viewInvalid, projInvalid bool
	viewLock, projLock       sync.RWMutex
}

func (camera *BaseCamera) Position() Vec3 {
	camera.viewLock.RLock()
	defer camera.viewLock.RUnlock()
	return camera.position
}

func (camera *BaseCamera) SetPosition(pos Vec3) {
	camera.viewLock.Lock()
	defer camera.viewLock.Unlock()
	camera.position = pos
	camera.viewInvalid = true
}

func (camera *BaseCamera) LookAt() Vec3 {
	camera.viewLock.RLock()
	defer camera.viewLock.RUnlock()
	return camera.lookAt
}

func (camera *BaseCamera) SetLookAt(pos Vec3) {
	camera.viewLock.Lock()
	defer camera.viewLock.Unlock()
	camera.lookAt = pos
	camera.viewInvalid = true
}

func (camera *BaseCamera) Up() Vec3 {
	camera.viewLock.RLock()
	defer camera.viewLock.RUnlock()
	return camera.up
}

func (camera *BaseCamera) SetUp(pos Vec3) {
	camera.viewLock.Lock()
	defer camera.viewLock.Unlock()
	camera.up = pos
	camera.viewInvalid = true
}

func (camera *BaseCamera) View() Mat4x4 {
	camera.viewLock.Lock()
	defer camera.viewLock.Unlock()
	if camera.viewInvalid {
		camera.view = Mat4x4(mgl32.LookAtV(
			mgl32.Vec3(camera.position),
			mgl32.Vec3(camera.lookAt),
			mgl32.Vec3(camera.up),
		))
		camera.viewInvalid = false
	}
	return camera.view
}

func (camera *BaseCamera) Far() float32 {
	camera.projLock.RLock()
	defer camera.projLock.RUnlock()
	return camera.far
}

func (camera *BaseCamera) SetFar(far float32) {
	camera.projLock.Lock()
	defer camera.projLock.Unlock()
	camera.far = far
	camera.projInvalid = true
}

func (camera *BaseCamera) Near() float32 {
	camera.projLock.RLock()
	defer camera.projLock.RUnlock()
	return camera.near
}

func (camera *BaseCamera) SetNear(near float32) {
	camera.projLock.Lock()
	defer camera.projLock.Unlock()
	camera.near = near
	camera.projInvalid = true
}

type OrhtoCamera struct {
	BaseCamera
	left, right, top, bottom, near, far float32
}

func (camera *OrhtoCamera) Viewport() (left, right, top, bottom float32) {
	camera.projLock.RLock()
	defer camera.projLock.RUnlock()
	return camera.left, camera.right, camera.top, camera.bottom
}

func (camera *OrhtoCamera) SetViewport(left, top, right, bottom float32) {
	camera.projLock.Lock()
	defer camera.projLock.Unlock()
	camera.left = left
	camera.right = right
	camera.top = top
	camera.bottom = bottom
	camera.projInvalid = true
}

func (camera *OrhtoCamera) Projection() Mat4x4 {
	camera.projLock.Lock()
	defer camera.projLock.Unlock()
	if camera.projInvalid {
		camera.proj = Mat4x4(mgl32.Ortho(
			camera.left,
			camera.right,
			camera.bottom,
			camera.top,
			camera.near,
			camera.far,
		))
		camera.projInvalid = false
	}
	return camera.proj
}
