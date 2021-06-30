package component

import "github.com/dhindustries/graal"

type Textured struct {
	texture graal.Texture
}

func (comp *Textured) Texture() graal.Texture {
	return comp.texture
}

func (comp *Textured) SetTexture(texture graal.Texture) {
	graal.Release(comp.texture)
	graal.Acquire(texture)
	comp.texture = texture
}
