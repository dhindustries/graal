package graal

import "github.com/go-gl/mathgl/mgl32"

type Vertex struct {
	Position  mgl32.Vec3
	Normal    mgl32.Vec3
	TexCoords mgl32.Vec2
	Color     mgl32.Vec4
}
