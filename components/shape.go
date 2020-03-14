package components

import "github.com/dhindustries/graal"

type Shaped struct {
	shape graal.Shape
}

func (component *Shaped) Shape() graal.Shape {
	return component.shape
}

func (component *Shaped) SetShape(shape graal.Shape) {
	if handle, ok := component.shape.(graal.Handle); ok {
		handle.Release()
	}
	if handle, ok := shape.(graal.Handle); ok {
		handle.Acquire()
	}
	component.shape = shape
}

func (component *Shaped) Dispose() {
	component.SetShape(nil)
}
