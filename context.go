package graal

import (
	"fmt"
	"runtime"

	"github.com/dhindustries/graal/queue"
)

type context struct {
	api Api
	app interface{}
	wnd Window
	rnr *queue.Runner
}

type appInitializer interface {
	init(*context) error
}

type appSetuper interface {
	Setup() error
}

type appPreparer interface {
	Prepare() error
}

type appUpdater interface {
	Update(dt float32)
}

type appRenderer interface {
	Render()
}

type appDisposer interface {
	Dispose()
}

func (ctx *context) run() error {
	go ctx.startQueue()
	defer ctx.stopQueue()
	defer ctx.cleanup()

	if err := ctx.initSystem(); err != nil {
		return err
	}
	defer ctx.finitSystem()

	if err := ctx.initWindow(); err != nil {
		return err
	}
	defer ctx.finitWindow()

	if err := ctx.appInit(); err != nil {
		return err
	}

	if err := ctx.initInput(); err != nil {
		return err
	}
	defer ctx.finitInput()

	if err := ctx.initGraphics(); err != nil {
		return err
	}
	defer ctx.finitGraphics()

	if err := ctx.appSetup(); err != nil {
		return err
	}
	defer ctx.appDispose()

	if err := ctx.wnd.Open(); err != nil {
		return err
	}
	if err := ctx.appPrepare(); err != nil {
		return err
	}
	ctx.log("Application start...\n")
	for ctx.running() {
		ctx.render()
		ctx.update()
	}
	ctx.log("Application done\n")
	return nil
}

func (ctx *context) running() bool {
	return ctx.wnd.IsOpen()
}

func (ctx *context) update() {
	ctx.wnd.PullMessages()
	ctx.appUpdate(1.0 / 60.0)
	ctx.updateInput()
}

func (ctx *context) render() {
	ctx.appRender()
	if ctx.api.RenderCommit != nil {
		go ctx.api.RenderCommit(&ctx.api)
	}
}

func (ctx *context) initSystem() error {
	if ctx.api.InitSystem != nil {
		ctx.log("Initializing system...\n")
		return ctx.api.InitSystem(&ctx.api)
	}
	return fmt.Errorf("system initialization required")
}

func (ctx *context) finitSystem() {
	if ctx.api.FinitSystem != nil {
		ctx.log("Disposing system...\n")
		ctx.api.FinitSystem(&ctx.api)
	}
}

func (ctx *context) initGraphics() error {
	if ctx.api.InitGraphics != nil {
		ctx.log("Initializing graphics...\n")
		return ctx.api.InitGraphics(&ctx.api, ctx.wnd)
	}
	return nil
}

func (ctx *context) finitGraphics() {
	if ctx.api.FinitGraphics != nil {
		ctx.log("Disposing graphics...\n")
		ctx.api.FinitGraphics(&ctx.api)
	}
}

func (ctx *context) initInput() error {
	if ctx.api.InitInput != nil {
		ctx.log("Initializing input...\n")
		return ctx.api.InitInput(&ctx.api, ctx.wnd)
	}
	return nil
}

func (ctx *context) finitInput() {
	if ctx.api.FinitInput != nil {
		ctx.log("Disposing input...\n")
		ctx.api.FinitInput(&ctx.api)
	}
}

func (ctx *context) initWindow() error {
	if ctx.api.NewWindow != nil {
		wnd, err := ctx.api.NewWindow(&ctx.api)
		if err != nil {
			return err
		}
		ctx.wnd = wnd
		return nil
	}
	return fmt.Errorf("window initializer not exists")
}

func (ctx *context) finitWindow() {
	if ctx.wnd != nil {
		ctx.wnd.Dispose()
	}
}

func (ctx *context) updateInput() {
	if ctx.api.UpdateInput != nil {
		ctx.api.UpdateInput(&ctx.api)
	}
}

func (ctx *context) appInit() error {
	if a, ok := ctx.app.(appInitializer); ok {
		if err := a.init(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *context) appSetup() error {
	if a, ok := ctx.app.(appSetuper); ok {
		if err := a.Setup(); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *context) appPrepare() error {
	if a, ok := ctx.app.(appPreparer); ok {
		if err := a.Prepare(); err != nil {
			return err
		}
		ctx.cleanup()
	}
	return nil
}

func (ctx *context) appUpdate(dt float32) {
	if a, ok := ctx.app.(appUpdater); ok {
		a.Update(dt)
	}
}

func (ctx *context) appRender() {
	if a, ok := ctx.app.(appRenderer); ok {
		a.Render()
	}
}

func (ctx *context) appDispose() {
	if a, ok := ctx.app.(appDisposer); ok {
		a.Dispose()
	}
}

func (ctx *context) cleanup() {
	if ctx.api.Cleanup != nil {
		ctx.api.Cleanup(&ctx.api)
	}
}

func (ctx *context) startQueue() {
	runtime.LockOSThread()
	ctx.rnr.Start()
	runtime.UnlockOSThread()
}

func (ctx *context) stopQueue() {
	ctx.rnr.Stop()
}

func (ctx *context) log(f string, v ...interface{}) {
	if ctx.api.Logf != nil {
		ctx.api.Logf(&ctx.api, f, v...)
	}
}

/*

func (context *context) start() error {
	defer context.Engine.RenderQueue.Break()
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

	return nil
}

func (context *context) initResourceManager() error {
	return context.ResourceManager.Initialize(&context.Engine)
}

func (context *context) initWindow() error {
	return context.Engine.RenderQueue.TryExec(func() error {
		return context.Window.Open()
	})
}

func (context *context) initGraphics() error {
	return context.Engine.RenderQueue.TryExec(func() error {
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
			context.RenderQueue.Push(context.render)
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
	context.RenderQueue.Invoke(context.Window.PullMessages)
	if app, ok := context.app.(appUpdate); ok {
		app.Update(1.0 / 60.0)
	}
	if context.Input != nil {
		context.Input.Update()
	}
}
*/
