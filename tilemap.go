package graal

type baseTilemap interface {
	Size() (w, h uint)
	SetSize(w, h uint)
	Tileset() Tileset
	SetTileset(v Tileset)
	SetTile(x, y, id uint)
	Tile(x, y uint) uint
	Mesh() Mesh
	TileParams(id uint) *ParamsWriter
}

type Tilemap interface {
	Handle
	baseTilemap
}

type TilemapResource interface {
	Resource
	baseTilemap
}
