package graal

type Image interface {
	Handle
	Size() (uint, uint)
	Get(x, y uint) Color
	Set(x, y uint, c Color)
	Data() []Color
}

type ImageResource interface {
	Resource
	Size() (uint, uint)
	Get(x, y uint) Color
	Set(x, y uint, c Color)
	Data() []Color
}
