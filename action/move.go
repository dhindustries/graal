package action

import (
	"time"

	"github.com/dhindustries/graal/utils"

	"github.com/go-gl/mathgl/mgl64"
)

type baseMove struct {
	t        time.Duration
	velocity mgl64.Vec3
}

type MoveBy struct {
	Position mgl64.Vec3
	Duration time.Duration
	baseMove
}

type MoveTo struct {
	Position mgl64.Vec3
	Duration time.Duration
	baseMove
}

type lookGetter interface {
	LookAt() mgl64.Vec3
}

type lookSetter interface {
	SetLookAt(mgl64.Vec3)
}

func (action *MoveTo) Run(t interface{}, dt time.Duration) bool {
	if action.t <= 0 {
		action.t = action.Duration
		pos, _ := utils.Position(t)
		action.velocity = action.Position.Sub(pos).Mul(1.0 / action.Duration.Seconds())
	}
	return action.baseMove.Run(t, dt)
}

func (action *MoveBy) Run(t interface{}, dt time.Duration) bool {
	if action.t <= 0 {
		action.t = action.Duration
		action.velocity = action.Position.Mul(1.0 / action.Duration.Seconds())
	}
	return action.baseMove.Run(t, dt)
}

func (action *baseMove) Run(t interface{}, dt time.Duration) bool {
	action.t -= dt
	if action.t > 0 {
		utils.AddPosition(t, action.velocity.Mul(dt.Seconds()))
		return false
	}
	return true
}

type Follow struct {
	Target interface{}
	Offset mgl64.Vec3
	init   bool
}

func (action *Follow) Run(t interface{}, dt time.Duration) bool {
	pos, _ := utils.Position(action.Target)
	mv := pos.Add(action.Offset)
	cur, _ := utils.Position(t)
	utils.AddPosition(t, mv)
	utils.AddLookAt(t, pos.Sub(cur))

	return false
}
