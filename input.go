package graal

type Input interface {
	Initialize(engine *Engine) error
	Keyboard() (Keyboard, error)
	Dispose()
	Update()
}
