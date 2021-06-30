package graal

type ShaderType = string

const (
	ShaderTypeVertex   = ShaderType("vertex")
	ShaderTypeFragment = ShaderType("fragment")
)

type baseShader interface {
	Type() ShaderType
}

type baseVertexShader interface {
}

type baseFragmentShader interface {
}

type Shader interface {
	Handle
	baseShader
}

type ShaderResource interface {
	Resource
	baseShader
}

type VertexShader interface {
	Shader
	baseVertexShader
}

type FragmentShader interface {
	Shader
	baseFragmentShader
}

type VertexShaderResource interface {
	ShaderResource
	baseVertexShader
}

type FragmentShaderResource interface {
	ShaderResource
	baseFragmentShader
}
