package graal

type baseMesh interface {
	SetVertexes(vx []Vertex)
	Build() error
}

type Mesh interface {
	Handle
	baseMesh
}

type MeshResource interface {
	Resource
	baseMesh
}
