package components

import (
	"github.com/dhindustries/graal"
)

type Meshed struct {
	mesh graal.Mesh
}

func (component *Meshed) Mesh() graal.Mesh {
	return component.mesh
}

func (component *Meshed) SetMesh(mesh graal.Mesh) {
	if component.mesh != nil {
		component.mesh.Release()
	}
	if mesh != nil {
		mesh.Acquire()
	}
	component.mesh = mesh
}

func (component *Meshed) Dispose() {
	component.SetMesh(nil)
}
