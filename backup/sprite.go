package graal

type baseSprite interface {
	Mesh() Mesh
	SetTexture(texture Texture)
	Texture() Texture
	Tileset() Tileset
	SetTileset(tileset Tileset)
	Frame() uint
	SetFrame(frame uint)
}

type Sprite interface {
	Handle
	baseSprite
}
