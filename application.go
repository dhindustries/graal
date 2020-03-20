package graal

type WindowApplication struct {
	api *Api
	wnd Window
}

type Application struct {
	WindowApplication
	kb Keyboard
	rs Resources
	rn Renderer
	fc Factory
	ms Mouse
}

func (app *WindowApplication) init(ctx *context) error {
	app.api = &ctx.api
	app.wnd = ctx.wnd
	return nil
}

func (app *WindowApplication) Window() Window {
	return app.wnd
}

func (app *WindowApplication) Close() {
	app.wnd.Close()
}

func (app *WindowApplication) Log(v interface{}) {
	app.api.Logf(app.api, "%v\n", v)
}

func (app *WindowApplication) Logf(fmt string, v ...interface{}) {
	app.api.Logf(app.api, "%v\n", v)
}

func (app *Application) Keyboard() Keyboard {
	if app.kb == nil {
		app.kb = &apiKeyboard{app.api}
	}
	return app.kb
}

func (app *Application) Resources() Resources {
	if app.rs == nil {
		app.rs = &apiResources{app.api}
	}
	return app.rs
}

func (app *Application) Renderer() Renderer {
	if app.rn == nil {
		app.rn = &apiRenderer{app.api}
	}
	return app.rn
}

func (app *Application) Factory() Factory {
	if app.fc == nil {
		app.fc = &apiFactory{app.api}
	}
	return app.fc
}

func (app *Application) Mouse() Mouse {
	if app.ms == nil {
		app.ms = &apiMouse{app.api}
	}
	return app.ms
}
