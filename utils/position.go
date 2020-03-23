package utils

import "github.com/go-gl/mathgl/mgl64"

type posGetter interface {
	Position() mgl64.Vec3
}

type posSetter interface {
	SetPosition(mgl64.Vec3)
}

func Position(object interface{}) (mgl64.Vec3, bool) {
	if v, ok := object.(posGetter); ok {
		return v.Position(), true
	}
	return mgl64.Vec3{}, false
}

func SetPosition(object interface{}, value mgl64.Vec3) bool {
	if v, ok := object.(posSetter); ok {
		v.SetPosition(value)
		return true
	}
	return false
}

func AddPosition(object interface{}, value mgl64.Vec3) bool {
	if pos, ok := Position(object); ok {
		return SetPosition(object, pos.Add(value))
	}
	return false
}
