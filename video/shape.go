package video

import (
	"github.com/dhindustries/graal"
	"github.com/dhindustries/graal/components"
	"github.com/go-gl/mathgl/mgl64"
)

type shape struct {
	graal.Handle
	api     *graal.Api
	mesh    graal.Mesh
	rebuild bool
	components.Colored
	components.Textured
}

func (sh *shape) Dispose() {
	sh.Textured.Dispose()
	sh.mesh.Release()
}

func (sh *shape) build(builder func() []graal.Vertex) {
	if sh.rebuild {
		sh.mesh.SetVertexes(builder())
		sh.rebuild = false
	}
}

type quad struct {
	shape
	left, right, top, bottom float64
}

func (sh *quad) Size() (width, height float64) {
	return sh.right - sh.left, sh.bottom - sh.top
}

func (sh *quad) SetSize(width, height float64) {
	sh.right, sh.bottom = sh.left+width, sh.top+height
	sh.rebuild = true
}

func (sh *quad) TopLeft() mgl64.Vec2 {
	return mgl64.Vec2{sh.left, sh.top}
}

func (sh *quad) SetTopLeft(value mgl64.Vec2) {
	sh.left, sh.top = value[0], value[1]
	sh.rebuild = true
}

func (sh *quad) BottomRight() mgl64.Vec2 {
	return mgl64.Vec2{sh.right, sh.bottom}
}

func (sh *quad) SetBottomRight(value mgl64.Vec2) {
	sh.right, sh.bottom = value[0], value[1]
	sh.rebuild = true
}

func (sh *quad) build() {
	sh.shape.build(func() []graal.Vertex {
		tl := graal.Vertex{
			Position:  mgl64.Vec3{sh.left, sh.top, 0},
			TexCoords: mgl64.Vec2{0, 0},
			Color:     mgl64.Vec4{1, 1, 1, 1},
		}
		tr := graal.Vertex{
			Position:  mgl64.Vec3{sh.right, sh.top, 0},
			TexCoords: mgl64.Vec2{1, 0},
			Color:     mgl64.Vec4{1, 1, 1, 1},
		}
		bl := graal.Vertex{
			Position:  mgl64.Vec3{sh.left, sh.bottom, 0},
			TexCoords: mgl64.Vec2{0, 1},
			Color:     mgl64.Vec4{1, 1, 1, 1},
		}
		br := graal.Vertex{
			Position:  mgl64.Vec3{sh.right, sh.bottom, 0},
			TexCoords: mgl64.Vec2{1, 1},
			Color:     mgl64.Vec4{1, 1, 1, 1},
		}
		return []graal.Vertex{
			tr, bl, tl,
			br, bl, tr,
		}
	})
}

func (sh *quad) Mesh() graal.Mesh {
	sh.build()
	return sh.mesh
}

func newQuad(api *graal.Api, left, right, top, bottom float64) (graal.Shape, error) {
	mesh, err := api.NewMesh(api, nil)
	if err != nil {
		return nil, err
	}
	q := &quad{
		shape: shape{
			Handle:  api.NewHandle(api),
			rebuild: true,
			mesh:    mesh,
		},
		left:   left,
		right:  right,
		top:    top,
		bottom: bottom,
	}
	q.SetColor(graal.ColorWhite)
	api.Handle(api, q)
	return q, nil
}
