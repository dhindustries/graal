package opengl

import (
	"github.com/dhindustries/graal"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Provider struct{}

func (Provider) Provide(api *graal.Api) error {
	renderer := renderer{}
	log := logger{api}

	prevInitSystem := api.InitSystem
	api.InitSystem = func(api *graal.Api) error {
		err := prevInitSystem(api)
		if err == nil {
			glfw.WindowHint(glfw.ContextVersionMajor, 4)
			glfw.WindowHint(glfw.ContextVersionMinor, 3)
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
		}

		return err
	}
	api.InitGraphics = func(api *graal.Api, wnd graal.Window) error {
		err := initGraphics(api, wnd)
		if err == nil {
			if log.enabled() {
				log.watch()
			}
			renderer.wnd = wnd
		}
		return err
	}
	api.FinitGraphics = func(api *graal.Api) {
		renderer.finit(api)
		finitGraphics(api)
	}

	// api.RenderEnqueue = func(api *graal.Api, o interface{}, t mgl32.Mat4) {
	// 	go renderer.enqueue(api, o, t)
	// }
	api.RenderEnqueue = renderer.enqueue
	api.RenderCommit = renderer.commit
	api.BindProgram = renderer.bindProgram
	api.SetCamera = renderer.setCamera

	api.NewShader = newShader
	api.NewProgram = newProgram
	api.NewMesh = newMesh
	api.NewTexture = newTexture

	api.SetResourceLoader(api, "texture/*", loadTextureResource)
	api.SetResourceLoader(api, "shader/vertex", loadVertexShaderResource)
	api.SetResourceLoader(api, "shader/fragment", loadFragmentShaderResource)
	api.SetResourceLoader(api, "program/*", loadProgramResource)

	return nil
}
