package graal

import "github.com/go-gl/mathgl/mgl64"

type Node interface {
	Handle
	ParentNode() Node
	Attach(v interface{})
	Detach(v interface{})
	List() []interface{}
	SetOrigin(origin mgl64.Vec3)
	SetPosition(position mgl64.Vec3)
	SetRotation(rotation mgl64.Vec3)
	SetScale(scale mgl64.Vec3)
	Origin() mgl64.Vec3
	Position() mgl64.Vec3
	Rotation() mgl64.Vec3
	Scale() mgl64.Vec3
	Transform() mgl64.Mat4
}
