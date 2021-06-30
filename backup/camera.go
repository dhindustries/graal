package graal

import "github.com/go-gl/mathgl/mgl64"

type Camera interface {
	Handle
	View() mgl64.Mat4
	Projection() mgl64.Mat4
	Position() mgl64.Vec3
	SetPosition(pos mgl64.Vec3)
	LookAt() mgl64.Vec3
	SetLookAt(pos mgl64.Vec3)
	Up() mgl64.Vec3
	SetUp(pos mgl64.Vec3)
	Far() float64
	SetFar(far float64)
	Near() float64
	SetNear(near float64)
}

type OrthoCamera interface {
	Camera
	Viewport() (left, right, top, bottom float64)
	SetViewport(left, top, right, bottom float64)
}
