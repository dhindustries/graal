package graal

import "log"

type Engine struct {
	Window          Window
	Graphics        Graphics
	Input           Input
	ResourceManager ResourceManager
	Logger          *log.Logger
	RenderQueue     Queue
}

func (engine Engine) Run(application interface{}) error {
	if engine.ResourceManager == nil {
		engine.ResourceManager = &BaseResourceManager{}
	}
	ctx := context{
		Engine: engine,
		app:    application,
	}
	return ctx.start()
}

func (engine Engine) log(v interface{}) {
	if engine.Logger != nil {
		engine.Logger.Println(v)
	}
}
