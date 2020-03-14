package graal

type Factory interface {
	Quad(left, top, right, bottom float32) Shape
}
