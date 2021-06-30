package graal

import "fmt"

type Resource interface {
	Handle
	Name() string
	Path() string
	Mime() Mime
}

type Resources interface {
	Load(mime Mime, name string) (Resource, error)
	LoadPrefab(name string) (Prefab, error)
	LoadFile(name string) (File, error)
	LoadImage(name string) (ImageResource, error)
	LoadTexture(name string) (TextureResource, error)
	LoadShader(t ShaderType, name string) (ShaderResource, error)
	LoadVertexShader(name string) (VertexShaderResource, error)
	LoadFragmentShader(name string) (FragmentShaderResource, error)
	LoadProgram(name string) (ProgramResource, error)
}

type apiResources struct {
	api *Api
}

func (res *apiResources) Load(mime Mime, name string) (Resource, error) {
	return res.api.LoadResource(res.api, mime, name)
}
func (res *apiResources) LoadPrefab(name string) (Prefab, error) {
	return res.api.LoadPrefab(res.api, name)
}

func (res *apiResources) LoadFile(name string) (File, error) {
	return res.api.LoadFile(res.api, name)
}

func (res *apiResources) LoadImage(name string) (ImageResource, error) {
	return res.api.LoadImage(res.api, name)
}

func (res *apiResources) LoadTexture(name string) (TextureResource, error) {
	return res.api.LoadTexture(res.api, name)
}

func (res *apiResources) LoadVertexShader(name string) (VertexShaderResource, error) {
	r, err := res.LoadShader(ShaderTypeVertex, name)
	if err != nil {
		return nil, err
	}
	if r.Type() != ShaderTypeVertex {
		r.Release()
		return nil, fmt.Errorf("shader %s is not a %s shader", r.Path(), ShaderTypeVertex)
	}
	return r, nil
}

func (res *apiResources) LoadFragmentShader(name string) (FragmentShaderResource, error) {
	r, err := res.LoadShader(ShaderTypeFragment, name)
	if err != nil {
		return nil, err
	}
	if r.Type() != ShaderTypeFragment {
		r.Release()
		return nil, fmt.Errorf("shader %s is not a %s shader", r.Path(), ShaderTypeFragment)
	}
	return r, nil
}

func (res *apiResources) LoadShader(t ShaderType, name string) (ShaderResource, error) {
	return res.api.LoadShader(res.api, t, name)
}

func (res *apiResources) LoadProgram(name string) (ProgramResource, error) {
	return res.api.LoadProgram(res.api, name)
}
