package gmath

import "github.com/go-gl/mathgl/mgl32"

type Mat4x4 mgl32.Mat4

func IdentityMat4x4() Mat4x4 {
	return Mat4x4(mgl32.Ident4())
}
