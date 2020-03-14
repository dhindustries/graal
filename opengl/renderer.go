package opengl

import (
	"fmt"
	"log"
	"runtime"

	"github.com/dhindustries/graal"
	"github.com/go-gl/gl/v4.3-compatibility/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type renderer struct {
	Window     graal.Window
	builder    builder
	camera     graal.Camera
	viewMatrix graal.Mat4x4
}

type textured interface {
	Texture() graal.Texture
}

type colored interface {
	Color() graal.Color
}

type transformed interface {
	Transformation() graal.Mat4x4
}

type shaped interface {
	Shape() graal.Shape
}

type clearableWindow interface {
	Clear(color graal.Color)
}

type swapableWindow interface {
	Swap()
}

func (*renderer) Dispose() {

}

func (renderer *renderer) Begin() {
	runtime.LockOSThread()
	if window, ok := renderer.Window.(clearableWindow); ok {
		window.Clear(graal.Color{})
	} else {
		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	}
	if renderer.camera != nil {
		renderer.viewMatrix = renderer.camera.View()
		proj := renderer.camera.Projection()
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadMatrixf(&proj[0])
	} else {
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		renderer.viewMatrix = graal.Mat4x4(mgl32.Ident4())
	}
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func (renderer *renderer) End() {
	if window, ok := renderer.Window.(swapableWindow); ok {
		window.Swap()
	}
	runtime.UnlockOSThread()
}

func (renderer *renderer) Render(object interface{}) {
	renderer.apply(object)
	switch obj := object.(type) {
	case shaped:
		object = obj.Shape()
	}
	renderer.apply(object)
	switch obj := object.(type) {
	case graal.Shape:
		renderer.renderShape(obj)
	}
}

func (renderer *renderer) Use(object interface{}) {
	if v, ok := object.(graal.Camera); ok {
		renderer.useCamera(v)
	}
}

func (renderer *renderer) useCamera(camera graal.Camera) {
	renderer.camera = camera
}

func (renderer *renderer) apply(object interface{}) {
	if v, ok := object.(transformed); ok {
		renderer.applyTransform(v)
	}
	if v, ok := object.(colored); ok {
		color := v.Color()
		gl.Color3f(color[0], color[1], color[2])
	}
	if v, ok := object.(textured); ok {
		renderer.applyTexture(v.Texture())
	}
}

func (renderer *renderer) renderShape(sh graal.Shape) {
	if shape, ok := sh.(*shape); ok {
		gl.Begin(shape.mode)
		for _, vx := range shape.points {
			gl.TexCoord2f(vx.TexCoords[0], vx.TexCoords[1])
			gl.Vertex4f(vx.Position[0], vx.Position[1], vx.Position[2], 1)
		}
		gl.End()
	} else {
		renderer.log(fmt.Sprintf("Only opengl shapes are supported"))
	}
}

func (*renderer) applyTransform(object transformed) {
	t := object.Transformation()
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadMatrixf(&t[0])
	// gl.MultMatrixf(&renderer.viewMatrix[0])
	// 	gl.LoadMatrixf(&renderer.viewMatrix[0])
	// 	gl.MultMatrixf(&t[0])
}

func (renderer *renderer) applyTexture(tex graal.Texture) {
	glid := uint32(0)
	if tex != nil {
		if tex, ok := tex.(*texture); ok {
			renderer.builder.buildTexture(tex)
			glid = tex.glid
		} else {
			renderer.log(fmt.Sprintf("Only opengl textures are supported"))
		}
	}

	gl.BindTexture(gl.TEXTURE_2D, glid)
}

func (renderer *renderer) log(v interface{}) {
	log.Println(v)
}
