package opengl

import (
	"github.com/dhindustries/graal"
	"github.com/dhindustries/graal/components"
)

type shape struct {
	mode   uint32
	points []graal.Vertex
	components.Colored
	components.Textured
}

func (shape *shape) Dispose() {
	shape.SetTexture(nil)
}
