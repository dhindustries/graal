package graal

type Shape interface {
	Disposable
	Texture() Texture
	SetTexture(texture Texture)
	Color() Color
	SetColor(color Color)
}
