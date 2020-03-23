package graal

import "github.com/go-gl/mathgl/mgl64"

type Vertex struct {
	Position  mgl64.Vec3
	Normal    mgl64.Vec3
	TexCoords mgl64.Vec2
	Color     mgl64.Vec4
}
