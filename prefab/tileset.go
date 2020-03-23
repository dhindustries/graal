package prefab

import (
	"encoding/xml"

	"github.com/dhindustries/graal"
)

type XMLTilesetPrefab struct {
	graal.Prefab
	api        *graal.Api
	XMLName    xml.Name `xml:"tileset"`
	Texture    string   `xml:"texture"`
	TileWidth  uint     `xml:"tile-width"`
	TileHeight uint     `xml:"tile-height"`
}

func (prefab *XMLTilesetPrefab) Spawn() (graal.Handle, error) {
	set, err := prefab.api.NewTileset(prefab.api)
	if err != nil {
		return nil, err
	}
	texPath := prefab.api.GetRelativePath(prefab.api, prefab, prefab.Texture)
	tex, err := prefab.api.LoadTexture(prefab.api, texPath)
	if err != nil {
		return nil, err
	}
	defer tex.Release()
	set.SetTexture(tex)
	set.SetTileSize(prefab.TileWidth, prefab.TileHeight)

	return set, nil
}

func loadXmlTilesetPrefab(api *graal.Api, prefab graal.Prefab) (graal.Prefab, error) {
	return &XMLTilesetPrefab{
		Prefab: prefab,
		api:    api,
	}, nil
}
