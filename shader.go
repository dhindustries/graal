package graal

type ShaderType = uint

const (
	ShaderTypeVertex   = ShaderType(1)
	ShaderTypeFragment = ShaderType(2)
)

type Shader interface {
	Type() ShaderType
	SetMatrix4x4(name string, matrix Mat4x4)
}
