package video

import (
	"github.com/dhindustries/graal"
	"github.com/dhindustries/graal/components"
	"github.com/dhindustries/graal/gmath"
	"github.com/go-gl/mathgl/mgl32"
)

type shape struct {
	graal.Handle
	components.Meshed
	components.Colored
	components.Textured
}

func (sh *shape) Dispose() {
	sh.Textured.Dispose()
	sh.Meshed.Dispose()
}

func newQuad(api *graal.Api, r gmath.Rect) (graal.Shape, error) {
	tl := graal.Vertex{
		Position:  mgl32.Vec3{r.Left, r.Top, 0},
		TexCoords: mgl32.Vec2{0, 0},
		Color:     mgl32.Vec4{1, 1, 1, 1},
	}
	tr := graal.Vertex{
		Position:  mgl32.Vec3{r.Right, r.Top, 0},
		TexCoords: mgl32.Vec2{1, 0},
		Color:     mgl32.Vec4{1, 1, 1, 1},
	}
	bl := graal.Vertex{
		Position:  mgl32.Vec3{r.Left, r.Bottom, 0},
		TexCoords: mgl32.Vec2{0, 1},
		Color:     mgl32.Vec4{1, 1, 1, 1},
	}
	br := graal.Vertex{
		Position:  mgl32.Vec3{r.Right, r.Bottom, 0},
		TexCoords: mgl32.Vec2{1, 1},
		Color:     mgl32.Vec4{1, 1, 1, 1},
	}
	mesh, err := api.NewMesh(api, []graal.Vertex{
		tr, bl, tl,
		br, bl, tr,
	})
	if err != nil {
		return nil, err
	}
	q := &shape{Handle: api.NewHandle(api)}
	q.SetMesh(mesh)
	api.Handle(api, q)
	return q, nil
}
