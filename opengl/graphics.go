package opengl

import (
	"fmt"
	"runtime"

	"github.com/dhindustries/graal"
	"github.com/go-gl/gl/v4.3-compatibility/gl"
)

type Graphics struct {
	initialized      bool
	window           graal.Window
	queue            graal.Queue
	rendererInstance *renderer
	factoryInstance  *factory
	builderInstance  *builder
	loaderInstance   *loader
}

func (graphics *Graphics) Initialize(engine *graal.Engine) error {
	if graphics.initialized {
		return fmt.Errorf("graphics device is already initialized")
	}
	runtime.LockOSThread()
	// defer runtime.UnlockOSThread()
	if err := gl.Init(); err != nil {
		return err
	}
	graphics.initialized = true
	graphics.window = engine.Window
	fmt.Println(engine.ResourceManager)
	if engine.ResourceManager != nil {
		fmt.Println("RM found")
		graphics.registerLoaders(engine.ResourceManager)
	}
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Enable(gl.TEXTURE_2D)
	return nil
}

func (graphics *Graphics) Dispose() {
	if !graphics.initialized {
		panic("graphics device is not initialized")
	}
	// runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	graphics.initialized = false
	graphics.window = nil
	gl.Finish()
}

func (graphics *Graphics) Renderer() (graal.Renderer, error) {
	if graphics.rendererInstance == nil {
		if graphics.window == nil {
			return nil, fmt.Errorf("window is not set")
		}
		builderInstance, err := graphics.builder()
		if err != nil {
			return nil, err
		}
		graphics.rendererInstance = &renderer{
			Window:  graphics.window,
			builder: new(builder),
		}
		*graphics.rendererInstance.builder = *builderInstance
		graphics.rendererInstance.builder.queue = nil
	}
	return graphics.rendererInstance, nil
}

func (graphics *Graphics) Factory() (graal.Factory, error) {
	if graphics.factoryInstance == nil {
		graphics.factoryInstance = &factory{}
	}
	return graphics.factoryInstance, nil
}

func (graphics *Graphics) registerLoaders(manager graal.ResourceManager) {
	loader, err := graphics.loader()
	if err == nil {
		manager.Register(graal.MimeTextureImage, loader.loadTexture)
	}
}

func (graphics *Graphics) loader() (*loader, error) {
	if graphics.loaderInstance == nil {
		builder, err := graphics.builder()
		if err != nil {
			return nil, err
		}
		graphics.loaderInstance = &loader{
			builder: builder,
			queue:   &graphics.queue,
		}
	}
	return graphics.loaderInstance, nil
}

func (graphics *Graphics) builder() (*builder, error) {
	if graphics.builderInstance == nil {
		graphics.builderInstance = &builder{
			queue: &graphics.queue,
		}
	}
	return graphics.builderInstance, nil
}
