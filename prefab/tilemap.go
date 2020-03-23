package prefab

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/dhindustries/graal"
)

type XMLCharMap struct {
	XMLName xml.Name
	Items   []struct {
		XMLName xml.Name
		Code    string `xml:"ch,attr"`
		Tile    uint   `xml:",chardata"`
	} `xml:"char"`
}

func (chmap *XMLCharMap) Get(code string) (uint, bool) {
	for _, i := range chmap.Items {
		if i.Code == code {
			return i.Tile, true
		}
	}
	return 0, false
}

type XMLTilesData struct {
	XMLName xml.Name
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

type XMLTileParam struct {
	ID     uint       `xml:"id,attr"`
	Params []XMLParam `xml:",any"`
}

type XMLTilemapPrefab struct {
	graal.Prefab
	api        *graal.Api
	XMLName    xml.Name       `xml:"tilemap"`
	Tileset    string         `xml:"tileset"`
	Width      uint           `xml:"width"`
	Height     uint           `xml:"height"`
	Map        XMLCharMap     `xml:"map"`
	Data       XMLTilesData   `xml:"tiles"`
	TileParams []XMLTileParam `xml:"tile-params>tile"`
}

func (prefab *XMLTilemapPrefab) Spawn() (graal.Handle, error) {
	tilesetPath := prefab.api.GetRelativePath(prefab.api, prefab, prefab.Tileset)
	tilesetPrefab, err := prefab.api.LoadPrefab(prefab.api, tilesetPath)
	if err != nil {
		return nil, err
	}
	defer tilesetPrefab.Release()

	tileset, err := tilesetPrefab.Spawn()
	if err != nil {
		return nil, err
	}
	defer tileset.Release()

	tilemap, err := prefab.api.NewTilemap(prefab.api)
	if err != nil {
		return nil, err
	}
	tilemap.SetTileset(tileset.(graal.Tileset))
	tilemap.SetSize(prefab.Width, prefab.Height)

	switch prefab.Data.Type {
	case "charmap":
		prefab.fillCharmap(tilemap)
	}

	if prefab.TileParams != nil {
		prefab.applyParams(tilemap)
	}

	return tilemap, nil
}

func (prefab *XMLTilemapPrefab) applyParams(tilemap graal.Tilemap) error {
	for _, tile := range prefab.TileParams {
		params := tilemap.TileParams(tile.ID)
		for _, param := range tile.Params {
			params.Set(param.Key, param.Value)
		}
	}
	return nil
}

func (prefab *XMLTilemapPrefab) fillCharmap(tilemap graal.Tilemap) error {
	content := strings.TrimSpace(prefab.Data.Content)
	if strings.HasPrefix(content, CDATABegin) {
		content = strings.TrimPrefix(content, CDATABegin)
		content = strings.TrimSuffix(content, CDATAEnd)
	}
	tiles := []uint{}
	for _, ch := range content {
		if tile, ok := prefab.Map.Get(string(ch)); ok {
			tiles = append(tiles, tile)
		}
	}
	if uint(len(tiles)) != prefab.Width*prefab.Height {
		return fmt.Errorf("incorrect data length")
	}
	for i, tile := range tiles {
		tilemap.SetTile(uint(i)%prefab.Width, uint(i)/prefab.Width, tile)
	}
	return nil
}

func loadXmlTilemapPrefab(api *graal.Api, prefab graal.Prefab) (graal.Prefab, error) {
	return &XMLTilemapPrefab{
		Prefab: prefab,
		api:    api,
	}, nil
}
