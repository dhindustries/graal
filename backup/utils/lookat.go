package utils

import "github.com/go-gl/mathgl/mgl64"

type lookatGetter interface {
	LookAt() mgl64.Vec3
}

type lookatSetter interface {
	SetLookAt(mgl64.Vec3)
}

func LookAt(object interface{}) (mgl64.Vec3, bool) {
	if v, ok := object.(lookatGetter); ok {
		return v.LookAt(), true
	}
	return mgl64.Vec3{}, false
}

func SetLookAt(object interface{}, value mgl64.Vec3) bool {
	if v, ok := object.(lookatSetter); ok {
		v.SetLookAt(value)
		return true
	}
	return false
}

func AddLookAt(object interface{}, value mgl64.Vec3) bool {
	if lookAt, ok := LookAt(object); ok {
		return SetLookAt(object, lookAt.Add(value))
	}
	return false
}
