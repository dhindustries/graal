package opengl

import (
	"github.com/dhindustries/graal"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

func vec2Trunc(v mgl64.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{float32(v[0]), float32(v[1])}
}

func vec3Trunc(v mgl64.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{float32(v[0]), float32(v[1]), float32(v[2])}
}

func vec4Trunc(v mgl64.Vec4) mgl32.Vec4 {
	return mgl32.Vec4{float32(v[0]), float32(v[1]), float32(v[2]), float32(v[3])}
}

func mat4Trunc(v mgl64.Mat4) mgl32.Mat4 {
	return mgl32.Mat4{
		float32(v[0]), float32(v[1]), float32(v[2]), float32(v[3]),
		float32(v[4]), float32(v[5]), float32(v[6]), float32(v[7]),
		float32(v[8]), float32(v[9]), float32(v[10]), float32(v[11]),
		float32(v[12]), float32(v[13]), float32(v[14]), float32(v[15]),
	}
}

func colorTrunc(v []graal.Color) []mgl32.Vec4 {
	res := make([]mgl32.Vec4, len(v))
	for i, e := range v {
		res[i] = vec4Trunc(mgl64.Vec4(e))
	}
	return res
}
