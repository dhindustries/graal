package video

import (
	"sync"

	"github.com/dhindustries/graal"
	"github.com/dhindustries/graal/components"
	"github.com/go-gl/mathgl/mgl64"
)

type sprite struct {
	graal.Handle
	components.Tiled
	components.Textured
	mesh      graal.Mesh
	rebuild   bool
	frame     uint
	texCoords [4]mgl64.Vec2
	l         sync.RWMutex
}

func (spr *sprite) Texture() graal.Texture {
	if ts := spr.Tileset(); ts != nil {
		if ts, ok := ts.(graal.MultiTileset); ok {
			return ts.TileTexture(spr.Frame())
		}
		return ts.Texture()
	}
	return spr.Textured.Texture()
}

func (spr *sprite) SetTileset(ts graal.Tileset) {
	spr.l.Lock()
	defer spr.l.Unlock()
	spr.Tiled.SetTileset(ts)
	spr.rebuild = true
}

func (spr *sprite) Frame() uint {
	spr.l.RLock()
	defer spr.l.RUnlock()
	return spr.frame
}

func (spr *sprite) SetFrame(frame uint) {
	spr.l.Lock()
	defer spr.l.Unlock()
	spr.frame = frame
	if ts := spr.Tileset(); ts != nil {
		tl, br := ts.GetTexCoords(spr.frame)
		spr.texCoords = [4]mgl64.Vec2{
			mgl64.Vec2{tl[0], tl[1]},
			mgl64.Vec2{br[0], tl[1]},
			mgl64.Vec2{tl[0], br[1]},
			mgl64.Vec2{br[0], br[1]},
		}
		spr.rebuild = true
	}
}

func (spr *sprite) Mesh() graal.Mesh {
	spr.l.Lock()
	defer spr.l.Unlock()
	spr.build()
	return spr.mesh
}

func (spr *sprite) build() {
	if spr.rebuild {
		tl := graal.Vertex{
			Position:  mgl64.Vec3{0, 0, 0},
			TexCoords: spr.texCoords[0],
			Color:     mgl64.Vec4{1, 1, 1, 1},
		}
		tr := graal.Vertex{
			Position:  mgl64.Vec3{1, 0, 0},
			TexCoords: spr.texCoords[1],
			Color:     mgl64.Vec4{1, 1, 1, 1},
		}
		bl := graal.Vertex{
			Position:  mgl64.Vec3{0, 1, 0},
			TexCoords: spr.texCoords[2],
			Color:     mgl64.Vec4{1, 1, 1, 1},
		}
		br := graal.Vertex{
			Position:  mgl64.Vec3{1, 1, 0},
			TexCoords: spr.texCoords[3],
			Color:     mgl64.Vec4{1, 1, 1, 1},
		}
		spr.mesh.SetVertexes([]graal.Vertex{
			tr, bl, tl,
			br, bl, tr,
		})
		spr.rebuild = false
	}
}

func (spr *sprite) Dispose() {
	spr.l.Lock()
	defer spr.l.Unlock()
	graal.Release(spr.Handle)
	graal.Release(spr.mesh)
	spr.SetTexture(nil)
	spr.SetTileset(nil)
	spr.mesh = nil
}
