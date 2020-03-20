package action

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type baseMove struct {
	t time.Duration
	s mgl32.Vec3
}

type MoveBy struct {
	Position mgl32.Vec3
	Duration time.Duration
	baseMove
}

type MoveTo struct {
	Position mgl32.Vec3
	Duration time.Duration
	baseMove
}

type posGetter interface {
	Position() mgl32.Vec3
}

type posSetter interface {
	SetPosition(mgl32.Vec3)
}

func (action *MoveTo) Run(t interface{}, dt time.Duration) bool {
	if action.t <= 0 {
		action.t = action.Duration
		var o mgl32.Vec3
		if v, ok := t.(posGetter); ok {
			o = v.Position()
		}
		action.s = action.Position.Sub(o).Mul(1.0 / float32(action.Duration.Seconds()))
	}
	return action.baseMove.Run(t, dt)
}

func (action *MoveBy) Run(t interface{}, dt time.Duration) bool {
	if action.t <= 0 {
		action.t = action.Duration
		action.s = action.Position.Mul(1.0 / float32(action.Duration.Seconds()))
	}
	return action.baseMove.Run(t, dt)
}

func (action *baseMove) Run(t interface{}, dt time.Duration) bool {
	action.t -= dt
	var p mgl32.Vec3
	if v, ok := t.(posGetter); ok {
		p = v.Position()
	}
	p = p.Add(action.s.Mul(float32(dt.Seconds())))
	if v, ok := t.(posSetter); ok {
		v.SetPosition(p)
	}
	return action.t <= 0
}
