package component

import "github.com/dhindustries/graal"

type Meshed struct {
	mesh graal.Mesh
}

func (comp *Meshed) Mesh() graal.Mesh {
	return comp.mesh
}

func (comp *Meshed) SetMesh(mesh graal.Mesh) {
	graal.Release(comp.mesh)
	graal.Acquire(mesh)
	comp.mesh = mesh
}

func (comp *Meshed) Dispose() {
	graal.Release(comp.mesh)
	comp.mesh = nil
}
