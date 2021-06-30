package opengl

import (
	"github.com/dhindustries/graal"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Library struct {
	renderer renderer
}

func (*Library) Name() string {
	return "opengl v4.6"
}

func (lib *Library) Install(proto *graal.ApiPrototype) error {

	proto.InitRenderer = lib.initRenderer
	proto.FinishRenderer = lib.finishRenderer

	proto.BeginRender = lib.renderer.begin
	proto.EndRender = lib.renderer.end
	proto.Render = lib.renderer.render
	proto.SetClearColor = lib.renderer.setClearColor
	proto.UseProgram = lib.renderer.useProgram
	proto.UseCamera = lib.renderer.useCamera

	proto.NewMesh = newMesh
	proto.NewTexture = newTexture
	proto.NewVertexShader = newVertexShader
	proto.NewFragmentShader = newFragmentShader
	proto.NewProgram = newProgram
	proto.NewText = newText

	return nil
}

func (*Library) Init(api graal.Api) error {
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.RedBits, 8)
	glfw.WindowHint(glfw.GreenBits, 8)
	glfw.WindowHint(glfw.BlueBits, 8)
	glfw.WindowHint(glfw.AlphaBits, 8)
	glfw.WindowHint(glfw.DepthBits, 32)
	glfw.WindowHint(glfw.StencilBits, 32)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.RefreshRate, 0)
	return nil
}

func (*Library) Finish(api graal.Api) {
}

func (*Library) initRenderer(api graal.Api) error {
	return api.TryInvoke(func(api graal.Api) error {
		if err := gl.Init(); err != nil {
			return err
		}
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
		gl.DebugMessageCallback(logGlMessage, nil)
		gl.DebugMessageControl(gl.DONT_CARE, gl.DONT_CARE, gl.DONT_CARE, 0, nil, true)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.Enable(gl.TEXTURE_2D)
		gl.Enable(gl.DEPTH_TEST)
		gl.Enable(gl.BLEND)
		gl.DepthFunc(gl.LESS)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		return nil
	})
}

func (lib *Library) finishRenderer(api graal.Api) {
	api.Release(lib.renderer)
	api.Invoke(func(api graal.Api) {
		gl.Finish()
	})
}
