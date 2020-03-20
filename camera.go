package graal

import "github.com/go-gl/mathgl/mgl32"

type Camera interface {
	Handle
	View() mgl32.Mat4
	Projection() mgl32.Mat4
	Position() mgl32.Vec3
	SetPosition(pos mgl32.Vec3)
	LookAt() mgl32.Vec3
	SetLookAt(pos mgl32.Vec3)
	Up() mgl32.Vec3
	SetUp(pos mgl32.Vec3)
	Far() float32
	SetFar(far float32)
	Near() float32
	SetNear(near float32)
}

type OrthoCamera interface {
	Camera
	Viewport() (left, right, top, bottom float32)
	SetViewport(left, top, right, bottom float32)
}
