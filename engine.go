package graal

import "log"

type Engine struct {
	Window          Window
	Graphics        Graphics
	Input           Input
	ResourceManager ResourceManager
	Logger          *log.Logger
}

func (engine Engine) Run(application interface{}) error {
	if engine.ResourceManager == nil {
		engine.ResourceManager = &BaseResourceManager{}
	}

	engine.log("Booting engine...")

	if engine.ResourceManager != nil {
		engine.log("Initializing resource manager...")
		if err := engine.ResourceManager.Initialize(&engine); err != nil {
			engine.log(err)
			return err
		}
	} else {
		engine.log("No resource manager")
	}

	engine.log("Initializing window...")
	if err := engine.Window.Open(); err != nil {
		engine.log(err)
		return err
	}
	defer engine.Window.Dispose()
	defer engine.Window.Close()

	engine.log("Initializing graphics...")
	if err := engine.Graphics.Initialize(&engine); err != nil {
		engine.log(err)
		return err
	}
	defer engine.Graphics.Dispose()

	renderer, err := engine.Graphics.Renderer()
	if err != nil {
		engine.log(err)
		return err
	}
	defer renderer.Dispose()

	if engine.Input != nil {
		engine.log("Initializing input...")
		if err := engine.Input.Initialize(&engine); err != nil {
			engine.log(err)
			return err
		}
		defer engine.Input.Dispose()
	} else {
		engine.log("No input")
	}

	if app, ok := application.(appSetup); ok {
		engine.log("Setup application...")
		app.Setup(engine)
	}

	if app, ok := application.(appPrepare); ok {
		engine.log("Preparing application...")
		app.Prepare()
		if engine.ResourceManager != nil {
			engine.ResourceManager.Cleanup()
		}
	}

	engine.log("Application start")
	for engine.Window.IsOpen() {
		engine.Window.PullMessages()
		renderer.Begin()
		if app, ok := application.(appRender); ok {
			app.Render()
		}
		renderer.End()
		if app, ok := application.(appUpdate); ok {
			app.Update(1.0 / 60.0)
		}
		if engine.Input != nil {
			engine.Input.Update()
		}
	}
	engine.log("Application done")

	if app, ok := application.(appDispose); ok {
		engine.log("Disposing application...")
		app.Dispose()
	}
	if engine.ResourceManager != nil {
		engine.ResourceManager.Cleanup()
	}

	return nil
}

func (engine Engine) log(v interface{}) {
	if engine.Logger != nil {
		engine.Logger.Println(v)
	}
}
