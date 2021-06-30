package graal

type Shape interface {
	Handle
	Texture() Texture
	SetTexture(texture Texture)
	Mesh() Mesh
	Color() Color
	SetColor(color Color)
}

type Quad interface {
	Shape
	Size() (width, height float64)
	SetSize(width, height float64)
}

type Circle interface {
	Shape
	Radius() float64
	SetRadius(readius float64)
	Angle() float64
	SetAngle(angle float64)
}
