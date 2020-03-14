package graal

import (
	"fmt"
)

type appSetup interface {
	Setup(engine Engine) error
}

type appDispose interface {
	Dispose()
}

type appUpdate interface {
	Update(dt float32)
}

type appRender interface {
	Render()
}

type appPrepare interface {
	Prepare()
}

type BaseApplication struct {
	engine Engine
}

func (app *BaseApplication) Setup(engine Engine) error {
	app.engine = engine
	return nil
}

func (app *BaseApplication) Resources() *Resources {
	if app.engine.ResourceManager == nil {
		app.Log("ResourceManager is not set")
		return nil
	}
	return &Resources{app.engine.ResourceManager}
}

func (app *BaseApplication) Factory() Factory {
	factory, err := app.engine.Graphics.Factory()
	if err != nil {
		app.Logf("Cannot access graphics factory: %s", err)
	}
	return factory
}

func (app *BaseApplication) Keyboard() Keyboard {
	if app.engine.Input != nil {
		keyboard, err := app.engine.Input.Keyboard()
		if err != nil {
			app.Log(fmt.Errorf("Cannot access keyboard: %s", err))
		}
		return keyboard
	}
	app.Log(fmt.Errorf("Input is not set"))
	return nil
}

func (app *BaseApplication) Window() Window {
	return app.engine.Window
}

func (app *BaseApplication) Renderer() Renderer {
	renderer, err := app.engine.Graphics.Renderer()
	if err != nil {
		app.Logf("Cannot access renderer: %s", err)
	}
	return renderer
}

func (app *BaseApplication) Close() {
	app.engine.Window.Close()
}

func (app *BaseApplication) Log(v interface{}) {
	if app.engine.Logger != nil {
		app.engine.Logger.Println(v)
	}
}

func (app *BaseApplication) Logf(fmt string, v ...interface{}) {
	if app.engine.Logger != nil {
		app.engine.Logger.Printf(fmt, v...)
	}
}
