package graal

type Renderer interface {
	Begin()
	End()
	Use(object interface{})
	Render(object interface{})
	Dispose()
}
