package video

import (
	"sync"

	"github.com/dhindustries/graal/memory"

	"github.com/dhindustries/graal"
	"github.com/go-gl/mathgl/mgl32"
)

type tileset struct {
	graal.Handle
	tex  graal.Texture
	w, h uint
	l    sync.RWMutex
}

func (ts *tileset) Size() (w, h uint) {
	ts.l.RLock()
	defer ts.l.RUnlock()
	if ts.tex != nil && ts.w > 0 && ts.h > 0 {
		tw, th := ts.tex.Size()
		return tw / ts.w, th / ts.h
	}
	return 0, 0
}

func (ts *tileset) TileSize() (w, h uint) {
	ts.l.RLock()
	defer ts.l.RUnlock()
	return ts.w, ts.h
}

func (ts *tileset) SetTileSize(w, h uint) {
	ts.l.Lock()
	defer ts.l.Unlock()
	ts.w = w
	ts.h = h
}

func (ts *tileset) Texture() graal.Texture {
	ts.l.RLock()
	defer ts.l.RUnlock()
	return ts.tex
}

func (ts *tileset) SetTexture(v graal.Texture) {
	ts.l.Lock()
	defer ts.l.Unlock()
	if ts.tex != nil {
		ts.tex.Release()
	}
	if v != nil {
		v.Acquire()
	}
	ts.tex = v
}

func (ts *tileset) GetTexCoords(id uint) (mgl32.Vec2, mgl32.Vec2) {
	w, h := ts.Size()
	if w == 0 || h == 0 {
		return mgl32.Vec2{}, mgl32.Vec2{1, 1}
	}
	if id >= w*h {
		id = 0
	}
	iw := 1.0 / float32(w)
	ih := 1.0 / float32(h)
	x := float32(id%w) * iw
	y := float32(id/w) * ih
	return mgl32.Vec2{x, y}, mgl32.Vec2{x + iw, y + ih}
}

func (ts *tileset) Dispose() {
	memory.Dispose(ts.Handle)
	ts.SetTexture(nil)
}

func newTileset(api *graal.Api) (graal.Tileset, error) {
	ts := &tileset{
		Handle: api.NewHandle(api),
		w:      1, h: 1,
	}
	api.Handle(api, ts)
	return ts, nil
}
