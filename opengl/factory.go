package opengl

import (
	"github.com/dhindustries/graal"
	"github.com/go-gl/gl/v4.3-compatibility/gl"
)

type factory struct{}

func (*factory) Quad(left, top, right, bottom float32) graal.Shape {
	s := &shape{
		mode: gl.QUADS,
		points: []graal.Vertex{
			graal.Vertex{graal.Vec3{left, top, 0}, graal.Vec2{0, 0}, graal.Vec3{}},
			graal.Vertex{graal.Vec3{left, bottom, 0}, graal.Vec2{0, 1}, graal.Vec3{}},
			graal.Vertex{graal.Vec3{right, bottom, 0}, graal.Vec2{1, 1}, graal.Vec3{}},
			graal.Vertex{graal.Vec3{right, top, 0}, graal.Vec2{1, 0}, graal.Vec3{}},
		},
	}
	s.SetColor(graal.Color{1, 1, 1, 1})
	return s
}
