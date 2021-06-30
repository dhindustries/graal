package graal

type Shader interface {
}

type VertexShader interface {
	Shader
}

type FragmentShader interface {
	Shader
}

type apiShader interface {
	NewVertexShader(content string) (VertexShader, error)
	NewFragmentShader(content string) (FragmentShader, error)
}

type protoShader struct {
	NewVertexShader   func(api Api, content string) (VertexShader, error)
	NewFragmentShader func(api Api, content string) (FragmentShader, error)
}

func NewVertexShader(content string) (VertexShader, error) {
	return api.NewVertexShader(content)
}

func (api *apiAdapter) NewVertexShader(content string) (VertexShader, error) {
	if api.proto.NewVertexShader == nil {
		panic("api.NewVertexShader is not implemented")
	}
	return api.proto.NewVertexShader(api, content)
}

func NewFragmentShader(content string) (FragmentShader, error) {
	return api.NewFragmentShader(content)
}

func (api *apiAdapter) NewFragmentShader(content string) (FragmentShader, error) {
	if api.proto.NewFragmentShader == nil {
		panic("api.NewFragmentShader is not implemented")
	}
	return api.proto.NewFragmentShader(api, content)
}
