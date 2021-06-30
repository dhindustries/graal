package graal

import (
	"sync"
	"sync/atomic"

	"github.com/go-gl/mathgl/mgl64"
)

type Camera interface {
	View() mgl64.Mat4
	Projection() mgl64.Mat4
	Position() mgl64.Vec3
	SetPosition(position mgl64.Vec3)
	LookAt() mgl64.Vec3
	SetLookAt(lookAt mgl64.Vec3)
	Up() mgl64.Vec3
	SetUp(up mgl64.Vec3)
	Far() float64
	SetFar(value float64)
	Near() float64
	SetNear(value float64)
}

type PerspectiveCamera interface {
	Camera
	AspectRatio() float64
	SetAspectRatio(value float64)
	FieldOfView() float64
	SetFieldOfView(value float64)
}

type OrthoCamera interface {
	Camera
	Viewport() (left, top, right, bottom float64)
	SetViewport(left, top, right, bottom float64)
}

type apiCamera interface {
	NewPerspectiveCamera() (PerspectiveCamera, error)
	NewOrthoCamera() (OrthoCamera, error)
	UseCamera(camera Camera)
}

type protoCamera struct {
	NewPerspectiveCamera func(api Api) (PerspectiveCamera, error)
	NewOrthoCamera       func(api Api) (OrthoCamera, error)
	UseCamera            func(api Api, camera Camera)
}

func NewPerspectiveCamera() (PerspectiveCamera, error) {
	return api.NewPerspectiveCamera()
}

func (api *apiAdapter) NewPerspectiveCamera() (PerspectiveCamera, error) {
	if api.proto.NewPerspectiveCamera == nil {
		panic("api.NewPerspectiveCamera is not implemented")
	}
	return api.proto.NewPerspectiveCamera(api)
}

func newPerspectiveCamera(api Api) (PerspectiveCamera, error) {
	return &perspectiveCamera{}, nil
}

func NewOrthoCamera() (OrthoCamera, error) {
	return api.NewOrthoCamera()
}

func (api *apiAdapter) NewOrthoCamera() (OrthoCamera, error) {
	if api.proto.NewOrthoCamera == nil {
		panic("api.NewOrthoCamera is not implemented")
	}
	return api.proto.NewOrthoCamera(api)
}

func newOrthoCamera(api Api) (OrthoCamera, error) {
	return &orthoCamera{}, nil
}

func UseCamera(camera Camera) {
	api.UseCamera(camera)
}

func (api *apiAdapter) UseCamera(camera Camera) {
	if api.proto.UseCamera == nil {
		panic("api.UseCamera is not implemented")
	}
	api.proto.UseCamera(api, camera)
}

type orthoCamera struct {
	camera
	viewport atomic.Value
}

func (cam *orthoCamera) Viewport() (left, right, top, bottom float64) {
	v := cam.viewport.Load().([4]float64)
	return v[0], v[1], v[2], v[3]
}

func (cam *orthoCamera) SetViewport(left, top, right, bottom float64) {
	cam.projMutex.RLock()
	defer cam.projMutex.RUnlock()
	cam.viewport.Store([4]float64{left, top, right, bottom})
	cam.projValid = false
}

func (cam *orthoCamera) Projection() mgl64.Mat4 {
	cam.projMutex.Lock()
	defer cam.projMutex.Unlock()
	if !cam.projValid {
		v := cam.viewport.Load().([4]float64)
		cam.proj = mgl64.Ortho(
			v[0], v[2], v[3], v[1],
			cam.near.Load().(float64),
			cam.far.Load().(float64),
		)
		cam.projValid = true
	}
	return cam.proj
}

type perspectiveCamera struct {
	camera
	aspect, fov atomic.Value
}

func (cam *perspectiveCamera) AspectRatio() float64 {
	return cam.aspect.Load().(float64)
}

func (cam *perspectiveCamera) SetAspectRatio(value float64) {
	cam.projMutex.RLock()
	defer cam.projMutex.RUnlock()
	cam.aspect.Store(value)
	cam.projValid = false
}

func (cam *perspectiveCamera) FieldOfView() float64 {
	return cam.fov.Load().(float64)
}

func (cam *perspectiveCamera) SetFieldOfView(value float64) {
	cam.projMutex.RLock()
	defer cam.projMutex.RUnlock()
	cam.fov.Store(value)
	cam.projValid = false
}

func (cam *perspectiveCamera) Projection() mgl64.Mat4 {
	cam.projMutex.Lock()
	defer cam.projMutex.Unlock()
	if !cam.projValid {
		cam.proj = mgl64.Perspective(
			cam.fov.Load().(float64),
			cam.aspect.Load().(float64),
			cam.near.Load().(float64),
			cam.far.Load().(float64),
		)
		cam.projValid = true
	}
	return cam.proj
}

type camera struct {
	view, proj           mgl64.Mat4
	viewValid, projValid bool
	pos, lookAt, up      atomic.Value
	near, far            atomic.Value
	viewMutex, projMutex sync.RWMutex
}

func (cam *camera) Position() mgl64.Vec3 {
	return cam.pos.Load().(mgl64.Vec3)
}

func (cam *camera) SetPosition(value mgl64.Vec3) {
	cam.viewMutex.RLock()
	defer cam.viewMutex.RUnlock()
	cam.pos.Store(value)
	cam.viewValid = false
}

func (cam *camera) LookAt() mgl64.Vec3 {
	return cam.lookAt.Load().(mgl64.Vec3)
}

func (cam *camera) SetLookAt(value mgl64.Vec3) {
	cam.viewMutex.RLock()
	defer cam.viewMutex.RUnlock()
	cam.lookAt.Store(value)
	cam.viewValid = false
}

func (cam *camera) Up() mgl64.Vec3 {
	return cam.up.Load().(mgl64.Vec3)
}

func (cam *camera) SetUp(value mgl64.Vec3) {
	cam.viewMutex.RLock()
	defer cam.viewMutex.RUnlock()
	cam.up.Store(value)
	cam.viewValid = false
}

func (cam *camera) Near() float64 {
	return cam.near.Load().(float64)
}

func (cam *camera) SetNear(value float64) {
	cam.projMutex.RLock()
	defer cam.projMutex.RUnlock()
	cam.near.Store(value)
	cam.projValid = false
}

func (cam *camera) Far() float64 {
	return cam.far.Load().(float64)
}

func (cam *camera) SetFar(value float64) {
	cam.projMutex.RLock()
	defer cam.projMutex.RUnlock()
	cam.far.Store(value)
	cam.projValid = false
}

func (cam *camera) View() mgl64.Mat4 {
	cam.viewMutex.Lock()
	defer cam.viewMutex.Unlock()
	if !cam.viewValid {
		cam.view = mgl64.LookAtV(
			cam.pos.Load().(mgl64.Vec3),
			cam.lookAt.Load().(mgl64.Vec3),
			cam.up.Load().(mgl64.Vec3),
		)
		cam.viewValid = true
	}
	return cam.view
}
