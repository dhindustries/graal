package graal

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Renderer interface {
	BindProgram(p Program)
	SetCamera(c Camera)
	Render(v interface{})
}

type apiRenderer struct {
	api *Api
}

func (r *apiRenderer) SetCamera(cam Camera) {
	r.api.SetCamera(r.api, cam)
}

func (r *apiRenderer) BindProgram(prog Program) {
	r.api.BindProgram(r.api, prog)
}

func (r *apiRenderer) Render(v interface{}) {
	r.api.RenderEnqueue(r.api, v, mgl64.Ident4())
}
