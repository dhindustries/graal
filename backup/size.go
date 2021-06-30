package graal

import "github.com/go-gl/mathgl/mgl64"

type Size struct {
	Left, Top, Right, Bottom float64
}

func (size *Size) Width() float64 {
	return size.Right - size.Bottom
}

func (size *Size) Height() float64 {
	return size.Bottom - size.Top;
}

func (size *Size) TopLeft() mgl64.Vec2 {
	return mgl64.Vec2{size.Left, size.Top}
}

func (size *Size) TopRight() mgl64.Vec2 {
	return mgl64.Vec2{size.Right, size.Top}
}

func (size *Size) BottomLeft() mgl64.Vec2 {
	return mgl64.Vec2{size.Left, size.Bottom}
}

func (size *Size) BottomRight() mgl64.Vec2 {
	return mgl64.Vec2{size.Right, size.Bottom}
}
