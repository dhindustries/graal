package graal

import (
	"fmt"
	"runtime"
	"time"

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
	Update(dt time.Duration)
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
	ctx.appUpdate(time.Second / 60)
	ctx.updateInput()
}

func (ctx *context) render() {
	ctx.appRender()
	if ctx.api.RenderCommit != nil {
		ctx.api.RenderCommit(&ctx.api)
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

func (ctx *context) appUpdate(dt time.Duration) {
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
