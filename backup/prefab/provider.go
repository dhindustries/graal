package prefab

import (
	"github.com/dhindustries/graal"
)

type Provider struct{}

func (Provider) Provide(api *graal.Api) error {
	api.Logf(api, "Installing prefabs...\n")
	xml := &xmlResourceLoader{make(map[string]graal.PrefabLoader)}

	api.LoadPrefab = loadPrefab
	api.SetResourceLoader(api, "prefab/xml", xml.loadXmlPefab)
	api.SetPrefabLoader = func(api *graal.Api, t, n string, l graal.PrefabLoader) {
		switch t {
		case "xml":
			xml.setLoader(api, n, l)
		}
	}
	xml.setLoader(api, "tileset", loadXmlTilesetPrefab)
	xml.setLoader(api, "tilemap", loadXmlTilemapPrefab)
	return nil
}
