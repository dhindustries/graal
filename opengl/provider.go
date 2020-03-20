package opengl

import (
	"github.com/dhindustries/graal"
)

type Provider struct{}

func (Provider) Provide(api *graal.Api) error {
	renderer := renderer{}
	log := logger{api}

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
