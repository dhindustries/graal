package opengl

import (
	"github.com/dhindustries/graal"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl64"
)

type meshed interface {
	Mesh() graal.Mesh
}

type transformed interface {
	Transformation() mgl64.Mat4
}

type textured interface {
	Texture() graal.Texture
}

type renderable interface {
	Render(api graal.Api, transformation mgl64.Mat4)
}

type renderer struct {
	cam                 graal.Camera
	prog                *program
	clearColor          mgl64.Vec4
	modelTransformation mgl64.Mat4
}

func (renderer *renderer) Dispose(api graal.Api) {
	if renderer.cam != nil {
		api.Release(renderer.cam)
		renderer.cam = nil
	}
	if renderer.prog != nil {
		api.Release(renderer.prog)
		renderer.prog = nil
		glUseProgram(api, 0)
	}
}

func (renderer *renderer) begin(api graal.Api) {
	glClear(api, renderer.clearColor)
	renderer.applyCamera(api)
}

func (renderer *renderer) end(api graal.Api) {

}

func (renderer *renderer) render(api graal.Api, object interface{}, transformation mgl64.Mat4) {
	if v, ok := object.(transformed); ok {
		transformation = transformation.Mul4(v.Transformation())
		//fmt.Printf("transformed: %v\n", transformation)
	}
	if obj, ok := object.(renderable); ok {
		obj.Render(api, transformation)
	} else {
		var obj = object
		var texture graal.Texture
		if v, ok := object.(meshed); ok {
			obj = v.Mesh()
		}
		if v, ok := object.(textured); ok {
			texture = v.Texture()
		}
		renderer.applyTransform(api, transformation)
		renderer.applyTexture(api, texture)
		renderer.renderObject(api, obj)
	}
}

func (renderer *renderer) renderObject(api graal.Api, object interface{}) {
	switch obj := object.(type) {
	case *mesh:
		renderer.renderMesh(api, obj)
	case *text:
		renderer.renderText(api, obj)
	}
}

func (renderer *renderer) setClearColor(api graal.Api, color graal.Color) {
	renderer.clearColor = mgl64.Vec4(color)
}

func (renderer *renderer) useProgram(api graal.Api, prog graal.Program) {
	if renderer.prog != nil {
		api.Release(renderer.prog)
	}
	if p, ok := prog.(*program); ok {
		api.Acquire(p)
		renderer.prog = p
		glUseProgram(api, renderer.prog.id)
	} else {
		renderer.prog = nil
		glUseProgram(api, 0)
	}
	renderer.applyCamera(api)
}

func (renderer *renderer) useCamera(api graal.Api, cam graal.Camera) {
	if renderer.cam != nil {
		api.Release(renderer.cam)
	}
	if cam != nil {
		api.Acquire(cam)
	}
	renderer.cam = cam
	renderer.applyCamera(api)
}

func (renderer *renderer) applyCamera(api graal.Api) {
	var view, proj mgl64.Mat4
	if renderer.cam != nil {
		view = renderer.cam.View()
		proj = renderer.cam.Projection()
	} else {
		ident := mgl64.Ident4()
		view = ident
		proj = ident
	}
	if renderer.prog != nil {
		glSetMat4f(api, renderer.prog.id, viewName, mat4Trunc(view))
		glSetMat4f(api, renderer.prog.id, projectionName, mat4Trunc(proj))
	}
}

func (renderer *renderer) applyTransform(api graal.Api, transform mgl64.Mat4) {
	renderer.modelTransformation = transform
	if renderer.prog != nil {
		glSetMat4f(api, renderer.prog.id, worldName, mat4Trunc(transform))
	}
}

func (renderer *renderer) applyTexture(api graal.Api, tex graal.Texture) {
	if t, ok := tex.(*texture); ok {
		glBindTexture(api, t.id)
	} else {
		glBindTexture(api, 0)
	}
}

func glClear(api graal.Api, color mgl64.Vec4) {
	api.Invoke(func(api graal.Api) {
		gl.ClearColor(float32(color[0]), float32(color[1]), float32(color[2]), float32(color[3]))
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	})
}
