package components

import (
	"github.com/dhindustries/graal"
)

type Tiled struct {
	tileset graal.Tileset
}

func (component *Tiled) Tileset() graal.Tileset {
	return component.tileset
}

func (component *Tiled) SetTileset(tileset graal.Tileset) {
	if component.tileset != nil {
		component.tileset.Release()
	}
	if tileset != nil {
		tileset.Acquire()
	}
	component.tileset = tileset
}

func (component *Tiled) Dispose() {
	component.SetTileset(nil)
}
