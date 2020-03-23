package controller

import (
	"sync"
	"time"

	"github.com/dhindustries/graal"
	"github.com/dhindustries/graal/utils"
	"github.com/go-gl/mathgl/mgl64"
)

type camera interface {
	LookAt() mgl64.Vec3
	SetLookAt(v mgl64.Vec3)
}

type positioned interface {
	Position() mgl64.Vec3
	SetPosition(v mgl64.Vec3)
}

type DPad struct {
	k          graal.Keyboard
	l, r, u, d graal.Key
	x, y       mgl64.Vec3
	m          sync.Mutex
	s          float64
}

func NewDPad(k graal.Keyboard) *DPad {
	return &DPad{
		k: k,
		l: graal.KeyLeft, r: graal.KeyRight,
		u: graal.KeyUp, d: graal.KeyDown,
		x: mgl64.Vec3{1, 0, 0},
		y: mgl64.Vec3{0, 1, 0},
		s: 1,
	}
}

func (p *DPad) SetLeftKey(k graal.Key) {
	p.m.Lock()
	defer p.m.Unlock()
	p.l = k
}

func (p *DPad) SetRightKey(k graal.Key) {
	p.m.Lock()
	defer p.m.Unlock()
	p.r = k
}

func (p *DPad) SetUpKey(k graal.Key) {
	p.m.Lock()
	defer p.m.Unlock()
	p.u = k
}

func (p *DPad) SetDownKey(k graal.Key) {
	p.m.Lock()
	defer p.m.Unlock()
	p.d = k
}

func (p *DPad) SetAxis(x, y mgl64.Vec3) {
	p.m.Lock()
	defer p.m.Unlock()
	p.x = x
	p.y = y
}

func (p *DPad) SetSpeed(v float64) {
	p.m.Lock()
	defer p.m.Unlock()
	p.s = v
}

func (p *DPad) Update(target interface{}, dt time.Duration) {
	p.m.Lock()
	defer p.m.Unlock()
	x, y := float64(0), float64(0)
	if p.k.IsDown(p.u) {
		y = y - 1
	}
	if p.k.IsDown(p.d) {
		y = y + 1
	}
	if p.k.IsDown(p.l) {
		x = x - 1
	}
	if p.k.IsDown(p.r) {
		x = x + 1
	}
	f := dt.Seconds() * p.s
	dir := p.x.Mul(x * f).Add(p.y.Mul(y * f))
	utils.AddPosition(target, dir)
	utils.AddLookAt(target, dir)
}
