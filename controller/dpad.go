package controller

import (
	"sync"
	"time"

	"github.com/dhindustries/graal"
	"github.com/go-gl/mathgl/mgl32"
)

type camera interface {
	LookAt() mgl32.Vec3
	SetLookAt(v mgl32.Vec3)
}

type positioned interface {
	Position() mgl32.Vec3
	SetPosition(v mgl32.Vec3)
}

type DPad struct {
	k          graal.Keyboard
	l, r, u, d graal.Key
	x, y       mgl32.Vec3
	t          interface{}
	m          sync.Mutex
	b          chan bool
	s          float32
}

func NewDPad(k graal.Keyboard) *DPad {
	return &DPad{
		k: k,
		l: graal.KeyLeft, r: graal.KeyRight,
		u: graal.KeyUp, d: graal.KeyDown,
		x: mgl32.Vec3{1, 0, 0},
		y: mgl32.Vec3{0, 1, 0},
		s: 1,
	}
}

func (p *DPad) SetTarget(v interface{}) {
	p.m.Lock()
	defer p.m.Unlock()
	if x, ok := p.t.(graal.Handle); ok {
		x.Release()
	}
	if x, ok := v.(graal.Handle); ok {
		x.Acquire()
	}
	p.t = v
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

func (p *DPad) SetAxis(x, y mgl32.Vec3) {
	p.m.Lock()
	defer p.m.Unlock()
	p.x = x
	p.y = y
}

func (p *DPad) SetSpeed(v float32) {
	p.m.Lock()
	defer p.m.Unlock()
	p.s = v
}

func (pad *DPad) Start(freq float32) {
	pad.m.Lock()
	if pad.b != nil {
		pad.m.Unlock()
		return
	}
	pad.b = make(chan bool)
	pad.m.Unlock()
	t := time.NewTicker(time.Second / time.Duration(freq))
	p := time.Now()
loop:
	for {
		select {
		case <-pad.b:
			t.Stop()
			break loop
		case t := <-t.C:
			pad.Update(t.Sub(p))
			p = t
		}
	}
	pad.m.Lock()
	defer pad.m.Unlock()
	close(pad.b)
	t.Stop()
	pad.b = nil
}

func (p *DPad) Stop() {
	p.m.Lock()
	defer p.m.Unlock()
	if p.b != nil {
		p.b <- true
	}
}

func (p *DPad) Update(dt time.Duration) {
	p.m.Lock()
	defer p.m.Unlock()

	if t, ok := p.t.(positioned); ok {
		x, y := float32(0), float32(0)
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
		f := float32(dt.Seconds()) * p.s
		d := p.x.Mul(x * f).Add(p.y.Mul(y * f))
		t.SetPosition(
			t.Position().Add(d),
		)
		if t, ok := p.t.(camera); ok {
			t.SetLookAt(
				t.LookAt().Add(d),
			)
		}
	}
}

func (p *DPad) Dispose() {
	p.m.Lock()
	defer p.m.Unlock()
	if p.b != nil {
		p.b <- true
	}
	graal.Release(p.t)
	p.k = nil
	p.t = nil
}
