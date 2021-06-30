package graal

import "github.com/go-gl/mathgl/mgl64"

type Mesh interface {
}

type apiMesh interface {
	NewMesh(vertexes []Vertex, indices []uint) (Mesh, error)
	NewQuad(topLeft, rightTop, leftBottom, rightBottom Vertex) (Mesh, error)
}

type protoMesh struct {
	NewMesh func(api Api, vertexes []Vertex, indices []uint) (Mesh, error)
	NewQuad func(api Api, topLeft, rightTop, leftBottom, rightBottom Vertex) (Mesh, error)
}

func NewMesh(vertexes []Vertex, indices []uint) (Mesh, error) {
	return api.NewMesh(vertexes, indices)
}

func (api *apiAdapter) NewMesh(vertexes []Vertex, indices []uint) (Mesh, error) {
	if api.proto.NewMesh == nil {
		panic("api.NewMesh is not implemented")
	}
	return api.proto.NewMesh(api, vertexes, indices)
}

func NewSimpleQuad(left, top, right, bottom float64) (Mesh, error) {
	return NewQuad(
		Vertex{
			Position:  mgl64.Vec3{left, top, 0},
			TexCoords: mgl64.Vec2{0, 0},
			Color:     mgl64.Vec4{1, 1, 1, 1},
		},
		Vertex{
			Position:  mgl64.Vec3{right, top, 0},
			TexCoords: mgl64.Vec2{1, 0},
			Color:     mgl64.Vec4{1, 1, 1, 1},
		},
		Vertex{
			Position:  mgl64.Vec3{left, bottom, 0},
			TexCoords: mgl64.Vec2{0, 1},
			Color:     mgl64.Vec4{1, 1, 1, 1},
		},
		Vertex{
			Position:  mgl64.Vec3{right, bottom, 0},
			TexCoords: mgl64.Vec2{1, 1},
			Color:     mgl64.Vec4{1, 1, 1, 1},
		},
	)
}

func NewQuad(topLeft, rightTop, leftBottom, rightBottom Vertex) (Mesh, error) {
	return api.NewQuad(topLeft, rightTop, leftBottom, rightBottom)
}

func (api *apiAdapter) NewQuad(topLeft, rightTop, leftBottom, rightBottom Vertex) (Mesh, error) {
	if api.proto.NewQuad == nil {
		panic("api.NewQuad is not implemented")
	}
	return api.proto.NewQuad(api, topLeft, rightTop, leftBottom, rightBottom)
}

func newQuad(api Api, topLeft, rightTop, leftBottom, rightBottom Vertex) (Mesh, error) {
	return api.NewMesh(
		[]Vertex{
			topLeft,
			rightTop,
			leftBottom,
			rightBottom,
		},
		[]uint{
			1, 2, 0,
			3, 2, 1,
		},
	)
}
