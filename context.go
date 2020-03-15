package graal

import (
	"runtime"
	"sync"
)

type context struct {
	Engine
	app            interface{}
	renderer       Renderer
	running, ready sync.WaitGroup
}

func (context *context) start() error {
	context.ready.Add(1)
	go context.runQueue()
	go context.runApp()
	return context.main()
}

func (context *context) main() error {
	context.log("Booting engine...")

	if context.ResourceManager != nil {
		context.log("Initializing resource manager...")
		if err := context.initResourceManager(); err != nil {
			context.log(err)
			return err
		}
	} else {
		context.log("No resource manager")
	}

	context.log("Initializing window...")
	if err := context.initWindow(); err != nil {
		context.log(err)
		return err
	}
	defer context.Window.Dispose()
	defer context.Window.Close()

	context.log("Initializing graphics...")
	if err := context.initGraphics(); err != nil {
		context.log(err)
		return err
	}
	defer context.Graphics.Dispose()

	renderer, err := context.Graphics.Renderer()
	if err != nil {
		context.log(err)
		return err
	}
	context.renderer = renderer
	defer context.renderer.Dispose()

	if context.Input != nil {
		context.log("Initializing input...")
		if err := context.initInput(); err != nil {
			context.log(err)
			return err
		}
		defer context.Input.Dispose()
	} else {
		context.log("No input")
	}

	if app, ok := context.app.(appSetup); ok {
		context.log("Setup application...")
		app.Setup(context.Engine)
	}

	if app, ok := context.app.(appPrepare); ok {
		context.log("Preparing application...")
		app.Prepare()
		if context.ResourceManager != nil {
			context.ResourceManager.Cleanup()
		}
	}

	context.ready.Done()
	context.running.Wait()

	if app, ok := context.app.(appDispose); ok {
		context.log("Disposing application...")
		app.Dispose()
	}

	if context.ResourceManager != nil {
		context.ResourceManager.Cleanup()
	}
	context.Engine.RenderQueue.Break()

	return nil
}

func (context *context) initResourceManager() error {
	return context.ResourceManager.Initialize(&context.Engine)
}

func (context *context) initWindow() error {
	return context.Engine.RenderQueue.Exec(func() error {
		return context.Window.Open()
	})
}

func (context *context) initGraphics() error {
	return context.Engine.RenderQueue.Exec(func() error {
		return context.Graphics.Initialize(&context.Engine)
	})
}

func (context *context) initInput() error {
	return context.Input.Initialize(&context.Engine)
}

func (context *context) runQueue() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	context.log("Starting queue...")
	context.RenderQueue.Run()
	context.log("Queue done...")
}

func (context *context) runApp() {
	if context.app != nil {
		context.running.Add(1)
		defer context.running.Done()
		context.ready.Wait()

		context.log("Application start")
		for context.Window.IsOpen() {
			context.update()
			context.RenderQueue.Exec(func() error {
				context.render()
				return nil
			})
		}
		context.log("Application done")
	}
}

func (context *context) render() {
	context.renderer.Begin()
	if app, ok := context.app.(appRender); ok {
		app.Render()
	}
	context.renderer.End()
}

func (context *context) update() {
	context.Window.PullMessages()
	if app, ok := context.app.(appUpdate); ok {
		app.Update(1.0 / 60.0)
	}
	if context.Input != nil {
		context.Input.Update()
	}
}
