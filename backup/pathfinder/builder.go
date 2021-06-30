package pathfinder

import (
	"github.com/dhindustries/graal"
)

type TileInfo struct {
	Available bool
	Weight    float64
}

type TilemapBuilder struct {
	Tiles map[uint]TileInfo
}

func (builder *TilemapBuilder) Build(tilemap graal.Tilemap) *Grid {
	w, h := tilemap.Size()
	grid := NewGrid(0, 0, int(w)+1, int(h)+1)
	for y := uint(0); y < h; y++ {
		for x := uint(0); x < w; x++ {
			info := builder.getTileInfo(tilemap, tilemap.Tile(x, y))
			grid.SetAvailable(int(x), int(y), info.Available)
			grid.SetWeight(int(x), int(y), info.Weight)
		}
	}
	return grid
}

func (builder *TilemapBuilder) getTileInfo(tilemap graal.Tilemap, id uint) TileInfo {
	info, ok := builder.Tiles[id]
	if !ok {
		params := tilemap.TileParams(id)
		if v, ok := params.GetBool("obstacle"); ok {
			info.Available = !v
		} else {
			info.Available = false
		}
		if v, ok := params.GetFloat("cost"); ok {
			info.Weight = v
		} else {
			info.Weight = 1
		}
		builder.Tiles[id] = info
	}
	return info
}
