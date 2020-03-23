package graal

import "github.com/go-gl/mathgl/mgl64"

type baseTileset interface {
	Dimmensions() (w, h uint)
	TileSize() (w, h uint)
	SetTileSize(w, h uint)
	Texture() Texture
	SetTexture(v Texture)
	GetTexCoords(tileID uint) (leftTop, bottomRight mgl64.Vec2)
}

type Tileset interface {
	Handle
	baseTileset
}

type MultiTileset interface {
	Handle
	baseTileset
	TileTexture(tileID uint) Texture
	SetTileTexture(tileID uint, v Texture)
}
