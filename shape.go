package graal

type Shape interface {
	Handle
	Texture() Texture
	SetTexture(texture Texture)
	Mesh() Mesh
	// Color() Color
	// SetColor(color Color)
}
