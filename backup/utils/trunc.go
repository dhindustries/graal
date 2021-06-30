package utils

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

func TruncVec4(v mgl64.Vec4) mgl32.Vec4 {
	return mgl32.Vec4{
		float32(v[0]),
		float32(v[1]),
		float32(v[2]),
		float32(v[3]),
	}
}
func TruncVec3(v mgl64.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{
		float32(v[0]),
		float32(v[1]),
		float32(v[2]),
	}
}
func TruncVec2(v mgl64.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{
		float32(v[0]),
		float32(v[1]),
	}
}
func TruncMat4(v mgl64.Mat4) mgl32.Mat4 {
	return mgl32.Mat4{
		float32(v[0]),
		float32(v[1]),
		float32(v[2]),
		float32(v[3]),
		float32(v[4]),
		float32(v[5]),
		float32(v[6]),
		float32(v[7]),
		float32(v[8]),
		float32(v[9]),
		float32(v[10]),
		float32(v[11]),
		float32(v[12]),
		float32(v[13]),
		float32(v[14]),
		float32(v[15]),
	}
}
