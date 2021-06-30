package components

import "github.com/dhindustries/graal"

type Textured struct {
	texture graal.Texture
}

func (component *Textured) Texture() graal.Texture {
	return component.texture
}

func (component *Textured) SetTexture(texture graal.Texture) {
	if handle, ok := component.texture.(graal.Handle); ok {
		handle.Release()
	}
	if handle, ok := texture.(graal.Handle); ok {
		handle.Acquire()
	}
	component.texture = texture
}

func (texturable *Textured) Dispose() {
	texturable.SetTexture(nil)
}
