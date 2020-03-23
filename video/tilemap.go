package video

import (
	"sync"

	"github.com/dhindustries/graal/memory"
	"github.com/go-gl/mathgl/mgl64"

	"github.com/dhindustries/graal"
)

type tilemap struct {
	graal.Handle
	s    graal.Tileset
	m    graal.Mesh
	v    bool
	w, h uint
	d    []uint
	l    sync.RWMutex
	p    map[uint]*graal.ParamsWriter
}

func newTilemap(api *graal.Api) (graal.Tilemap, error) {
	m, err := api.NewMesh(api, nil)
	if err != nil {
		return nil, err
	}
	tm := &tilemap{
		Handle: api.NewHandle(api),
		m:      m,
	}
	api.Handle(api, tm)
	return tm, nil
}

func (tm *tilemap) Mesh() graal.Mesh {
	tm.l.Lock()
	defer tm.l.Unlock()
	if !tm.v {
		tm.m.SetVertexes(tm.buildVertexes())
		tm.v = true
	}
	return tm.m
}

func (tm *tilemap) Texture() graal.Texture {
	if s := tm.Tileset(); s != nil {
		return s.Texture()
	}
	return nil
}

func (tm *tilemap) Size() (w, h uint) {
	tm.l.RLock()
	defer tm.l.RUnlock()
	return tm.w, tm.h
}

func (tm *tilemap) SetSize(w, h uint) {
	tm.l.Lock()
	defer tm.l.Unlock()
	if tm.w != w || tm.h != h {
		tm.resize(w, h)
		tm.v = false
	}
}

func (tm *tilemap) Tileset() graal.Tileset {
	tm.l.RLock()
	defer tm.l.RUnlock()
	return tm.s
}

func (tm *tilemap) SetTileset(v graal.Tileset) {
	tm.l.Lock()
	defer tm.l.Unlock()
	if tm.s != nil {
		tm.s.Release()
	}
	if v != nil {
		v.Acquire()
	}
	tm.s = v
	tm.v = false
}

func (tm *tilemap) SetTile(x, y, id uint) {
	tm.l.Lock()
	defer tm.l.Unlock()
	if x < tm.w && y < tm.h {
		tm.d[tm.idx(x, y, tm.w, tm.h)] = id
		tm.v = false
	}
}

func (tm *tilemap) Tile(x, y uint) uint {
	tm.l.Lock()
	defer tm.l.Unlock()
	if x < tm.w && y < tm.h {
		return tm.d[tm.idx(x, y, tm.w, tm.h)]
	}
	return 0
}

func (tm *tilemap) TileParams(id uint) *graal.ParamsWriter {
	tm.l.Lock()
	defer tm.l.Unlock()
	if id >= uint(len(tm.d)) {
		return nil
	}
	if tm.p == nil {
		tm.p = make(map[uint]*graal.ParamsWriter)
	}
	params, ok := tm.p[id]
	if !ok {
		params = new(graal.ParamsWriter)
		tm.p[id] = params
	}
	return params
}

func (tm *tilemap) Dispose() {
	memory.Dispose(tm.Handle)
	tm.SetTileset(nil)
	tm.l.Lock()
	defer tm.l.Unlock()
	if tm.m != nil {
		tm.m.Release()
		tm.m = nil
	}
}

func (tm *tilemap) resize(w, h uint) {
	d := make([]uint, w*h)
	if tm.d != nil && tm.w != 0 && tm.h != 0 {
		cw := tm.min(w, tm.w)
		ch := tm.min(h, tm.h)
		for x := uint(0); x < cw; x++ {
			for y := uint(0); y < ch; y++ {
				d[tm.idx(x, y, w, h)] = tm.d[tm.idx(x, y, tm.w, tm.h)]
			}
		}
	}
	tm.d = d
	tm.w = w
	tm.h = h
}

func (tm *tilemap) buildVertexes() []graal.Vertex {
	list := make([]graal.Vertex, 0, tm.w*tm.h*6)
	if tm.s != nil {
		for y := uint(0); y < tm.h; y++ {
			for x := uint(0); x < tm.w; x++ {
				id := tm.d[tm.idx(x, y, tm.w, tm.h)]
				u, v := tm.s.GetTexCoords(id)

				tl := graal.Vertex{
					Position:  mgl64.Vec3{float64(x), float64(y), 0},
					TexCoords: u,
				}
				tr := graal.Vertex{
					Position:  mgl64.Vec3{float64(x + 1), float64(y), 0},
					TexCoords: mgl64.Vec2{v[0], u[1]},
				}
				bl := graal.Vertex{
					Position:  mgl64.Vec3{float64(x), float64(y + 1), 0},
					TexCoords: mgl64.Vec2{u[0], v[1]},
				}
				br := graal.Vertex{
					Position:  mgl64.Vec3{float64(x + 1), float64(y + 1), 0},
					TexCoords: v,
				}
				list = append(
					list,
					tr, bl, tl,
					br, bl, tr,
				)
			}
		}
	}
	return list
}

func (tm *tilemap) min(a, b uint) uint {
	if b > a {
		return a
	}
	return b
}

func (tm *tilemap) idx(x, y, w, h uint) uint {
	return y*w + x
}
