package graal

type baseSprite interface {
	Mesh() Mesh
	Texture() Texture
	SetTexture(texture Texture)
}

type Sprite interface {
	Handle
	baseSprite
}

type SpriteResource interface {
	Resource
	baseSprite
}
