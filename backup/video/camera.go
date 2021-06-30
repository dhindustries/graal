package video

import (
	"sync"

	"github.com/dhindustries/graal"

	"github.com/go-gl/mathgl/mgl64"
)

type baseCamera struct {
	graal.Handle
	position, lookAt, up     mgl64.Vec3
	far, near                float64
	view, proj               mgl64.Mat4
	viewInvalid, projInvalid bool
	viewLock, projLock       sync.RWMutex
}

type orthoCamera struct {
	baseCamera
	left, right, top, bottom float64
}

func newOrthoCamera(api *graal.Api) (graal.OrthoCamera, error) {
	return &orthoCamera{
		baseCamera: baseCamera{
			Handle: api.NewHandle(api),
		},
	}, nil
}

func (camera *orthoCamera) Viewport() (left, right, top, bottom float64) {
	camera.projLock.RLock()
	defer camera.projLock.RUnlock()
	return camera.left, camera.right, camera.top, camera.bottom
}

func (camera *orthoCamera) SetViewport(left, top, right, bottom float64) {
	camera.projLock.Lock()
	defer camera.projLock.Unlock()
	camera.left = left
	camera.right = right
	camera.top = top
	camera.bottom = bottom
	camera.projInvalid = true
}

func (camera *orthoCamera) Projection() mgl64.Mat4 {
	camera.projLock.Lock()
	defer camera.projLock.Unlock()
	if camera.projInvalid {
		camera.proj = mgl64.Ortho(
			camera.left,
			camera.right,
			camera.bottom,
			camera.top,
			camera.near,
			camera.far,
		)
		camera.projInvalid = false
	}
	return camera.proj
}

func (camera *baseCamera) Position() mgl64.Vec3 {
	camera.viewLock.RLock()
	defer camera.viewLock.RUnlock()
	return camera.position
}

func (camera *baseCamera) SetPosition(pos mgl64.Vec3) {
	camera.viewLock.Lock()
	defer camera.viewLock.Unlock()
	camera.position = pos
	camera.viewInvalid = true
}

func (camera *baseCamera) LookAt() mgl64.Vec3 {
	camera.viewLock.RLock()
	defer camera.viewLock.RUnlock()
	return camera.lookAt
}

func (camera *baseCamera) SetLookAt(pos mgl64.Vec3) {
	camera.viewLock.Lock()
	defer camera.viewLock.Unlock()
	camera.lookAt = pos
	camera.viewInvalid = true
}

func (camera *baseCamera) Up() mgl64.Vec3 {
	camera.viewLock.RLock()
	defer camera.viewLock.RUnlock()
	return camera.up
}

func (camera *baseCamera) SetUp(pos mgl64.Vec3) {
	camera.viewLock.Lock()
	defer camera.viewLock.Unlock()
	camera.up = pos
	camera.viewInvalid = true
}

func (camera *baseCamera) View() mgl64.Mat4 {
	camera.viewLock.Lock()
	defer camera.viewLock.Unlock()
	if camera.viewInvalid {
		camera.view = mgl64.Mat4(mgl64.LookAtV(
			mgl64.Vec3(camera.position),
			mgl64.Vec3(camera.lookAt),
			mgl64.Vec3(camera.up),
		))
		camera.viewInvalid = false
	}
	return camera.view
}

func (camera *baseCamera) Far() float64 {
	camera.projLock.RLock()
	defer camera.projLock.RUnlock()
	return camera.far
}

func (camera *baseCamera) SetFar(far float64) {
	camera.projLock.Lock()
	defer camera.projLock.Unlock()
	camera.far = far
	camera.projInvalid = true
}

func (camera *baseCamera) Near() float64 {
	camera.projLock.RLock()
	defer camera.projLock.RUnlock()
	return camera.near
}

func (camera *baseCamera) SetNear(near float64) {
	camera.projLock.Lock()
	defer camera.projLock.Unlock()
	camera.near = near
	camera.projInvalid = true
}
