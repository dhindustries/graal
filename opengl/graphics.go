package opengl

import (
	"fmt"
	"runtime"

	"github.com/dhindustries/graal"
	"github.com/go-gl/gl/v4.3-compatibility/gl"
)

type Graphics struct {
	initialized bool
	window      graal.Window
	renderer    *renderer
}

func (graphics *Graphics) Initialize(engine *graal.Engine) error {
	if graphics.initialized {
		return fmt.Errorf("graphics device is already initialized")
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := gl.Init(); err != nil {
		return err
	}
	graphics.initialized = true
	graphics.window = engine.Window
	if engine.ResourceManager != nil {
		engine.ResourceManager.Register(graal.MimeTextureImage, textureImageLoader)
	}
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Enable(gl.TEXTURE_2D)
	return nil
}

func (graphics *Graphics) Dispose() {
	if !graphics.initialized {
		panic("graphics device is not initialized")
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	graphics.initialized = false
	graphics.window = nil
	gl.Finish()
}

func (graphics *Graphics) Renderer() (graal.Renderer, error) {
	if graphics.renderer == nil {
		if graphics.window == nil {
			return nil, fmt.Errorf("window is not set")
		}
		graphics.renderer = &renderer{Window: graphics.window}
	}
	return graphics.renderer, nil
}

func (graphics *Graphics) Factory() (graal.Factory, error) {
	return &factory{}, nil
}
