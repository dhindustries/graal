package graal

import "github.com/go-gl/mathgl/mgl32"

type Node interface {
	Handle
	ParentNode() Node
	Attach(v interface{})
	Detach(v interface{})
	List() []interface{}
	SetOrigin(origin mgl32.Vec3)
	SetPosition(position mgl32.Vec3)
	SetRotation(rotation mgl32.Vec3)
	SetScale(scale mgl32.Vec3)
	Origin() mgl32.Vec3
	Position() mgl32.Vec3
	Rotation() mgl32.Vec3
	Scale() mgl32.Vec3
	Transform() mgl32.Mat4
}
