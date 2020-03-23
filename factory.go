package graal

type Factory interface {
	Node() (Node, error)
	Mesh(p []Vertex) (Mesh, error)
	Texture(w, h uint) (Texture, error)
	Quad(left, top, right, bottom float64) (Shape, error)
	Program() (Program, error)
	OrthoCamera() (OrthoCamera, error)
	Tileset() (Tileset, error)
	Tilemap() (Tilemap, error)
}

type apiFactory struct {
	api *Api
}

func (f *apiFactory) Quad(l, t, r, b float64) (Shape, error) {
	return f.api.NewQuad(f.api, l, r, t, b)
}

func (f *apiFactory) Program() (Program, error) {
	return f.api.NewProgram(f.api)
}

func (f *apiFactory) Texture(w, h uint) (Texture, error) {
	return f.api.NewTexture(f.api, w, h)
}

func (f *apiFactory) Node() (Node, error) {
	return f.api.NewNode(f.api)
}

func (f *apiFactory) Mesh(p []Vertex) (Mesh, error) {
	return f.api.NewMesh(f.api, p)
}

func (f *apiFactory) OrthoCamera() (OrthoCamera, error) {
	return f.api.NewOrthoCamera(f.api)
}
func (f *apiFactory) Tileset() (Tileset, error) {
	return f.api.NewTileset(f.api)
}
func (f *apiFactory) Tilemap() (Tilemap, error) {
	return f.api.NewTilemap(f.api)
}
