package graal

import (
	"github.com/go-gl/mathgl/mgl32"
)

type baseTileset interface {
	Size() (w, h uint)
	TileSize() (w, h uint)
	SetTileSize(w, h uint)
	Texture() Texture
	SetTexture(v Texture)
	GetTexCoords(tileID uint) (leftTop, bottomRight mgl32.Vec2)
}

type Tileset interface {
	Handle
	baseTileset
}

type TilesetResource interface {
	Resource
	baseTileset
}
