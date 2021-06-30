package graal

import (
	"github.com/go-gl/mathgl/mgl32"
)

type baseProgram interface {
	SetMat4f(name string, v mgl32.Mat4)
	SetVec4f(name string, v mgl32.Vec4)
	SetVertexShader(shader VertexShader)
	VertexShader() VertexShader
	SetFragmentShader(shader FragmentShader)
	FragmentShader() FragmentShader
	Compile() error
}

type Program interface {
	Handle
	baseProgram
}

type ProgramResource interface {
	Resource
	baseProgram
}
