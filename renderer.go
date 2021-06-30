package graal

import "github.com/go-gl/mathgl/mgl64"

type apiRenderer interface {
	initRenderer() error
	finishRenderer()
	beginRender()
	endRender()
	Render(object interface{}, transformation mgl64.Mat4)
	SetClearColor(color Color)
}

type protoRenderer struct {
	InitRenderer   func(api Api) error
	FinishRenderer func(api Api)
	BeginRender    func(api Api)
	EndRender      func(api Api)
	Render         func(api Api, object interface{}, transformation mgl64.Mat4)
	SetClearColor  func(api Api, color Color)
}

func SetClearColor(color Color) {
	api.SetClearColor(color)
}

func (api *apiAdapter) SetClearColor(color Color) {
	if api.proto.SetClearColor == nil {
		panic("api.SetClearColor is not implemented")
	}
	api.proto.SetClearColor(api, color)
}

func Render(object interface{}, transformation mgl64.Mat4) {
	api.Render(object, transformation)
}

func (api *apiAdapter) Render(object interface{}, transformation mgl64.Mat4) {
	if api.proto.Render == nil {
		panic("api.Render is not implemented")
	}
	api.proto.Render(api, object, transformation)
}

func (api *apiAdapter) initRenderer() error {
	//if api.proto.InitRenderer == nil {
	//	panic("api.InitRenderer is not implemented")
	//}
	//return api.proto.InitRenderer(api)
	if api.proto.InitRenderer != nil {
		return api.proto.InitRenderer(api)
	}
	return nil
}

func (api *apiAdapter) finishRenderer() {
	//if api.proto.FinishRenderer == nil {
	//	panic("api.FinishRenderer is not implemented")
	//}
	//api.proto.FinishRenderer(api)
	if api.proto.FinishRenderer != nil {
		api.proto.FinishRenderer(api)
	}
}

func (api *apiAdapter) beginRender() {
	if api.proto.BeginRender == nil {
		panic("api.BeginRender is not implemented")
	}
	api.proto.BeginRender(api)
}

func (api *apiAdapter) endRender() {
	if api.proto.EndRender == nil {
		panic("api.EndRender is not implemented")
	}
	api.proto.EndRender(api)
}
