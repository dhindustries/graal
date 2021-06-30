package graal

import "time"

var timeStep = time.Second / 60.0

type setuper interface {
	Setup() error
}

type loader interface {
	Load() error
}

type unloader interface {
	Unload()
}

type updater interface {
	Update(deltaTime time.Duration)
}

type renderer interface {
	Render()
}

func start(app interface{}) error {
	if v, ok := app.(setuper); ok {
		if err := v.Setup(); err != nil {
			return err
		}
	}
	defer Release(app)
	if v, ok := app.(loader); ok {
		if err := v.Load(); err != nil {
			return err
		}
	}
	if err := api.openWindow(800, 600, "graal"); err != nil {
		return err
	}
	defer api.closeWindow()
	if v, ok := app.(unloader); ok {
		defer v.Unload()
	}
	var clock Clock
	var acc time.Duration
	clock.Reset()
	for api.isWindowOpen() {
		api.updateWindow()
		if v, ok := app.(updater); ok {
			for acc += clock.Elapsed(); acc >= timeStep; acc -= timeStep {
				v.Update(timeStep)
			}
		}

		api.beginRender()
		if v, ok := app.(renderer); ok {
			v.Render()
		}
		api.endRender()

		api.renderWindow()
	}
	return nil
}
