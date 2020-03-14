package graal

type Graphics interface {
	Initialize(e *Engine) error
	Renderer() (Renderer, error)
	Factory() (Factory, error)
	Dispose()
}
