package prefab

import (
	"encoding/xml"
	"fmt"

	"github.com/dhindustries/graal"
)

const (
	CDATABegin = "<![CDATA["
	CDATAEnd   = "]]>"
)

type xmlResourceLoader struct {
	loaders map[string]graal.PrefabLoader
}

func (resLoader *xmlResourceLoader) setLoader(api *graal.Api, name string, loader graal.PrefabLoader) {
	resLoader.loaders[name] = loader
}

func (resLoader *xmlResourceLoader) loadXmlPefab(api *graal.Api, resource graal.Resource) (graal.Resource, error) {
	file, err := api.LoadFile(api, resource.Path())
	if err != nil {
		return nil, err
	}
	defer file.Release()
	decoder := xml.NewDecoder(file)
	i, err := decoder.Token()
	for ; err == nil; i, err = decoder.Token() {
		switch token := i.(type) {
		case xml.StartElement:
			if reader, ok := resLoader.loaders[token.Name.Local]; ok {
				prefab, err := reader(api, &basePrefab{Resource: resource})
				if err != nil {
					return nil, err
				}
				if err := decoder.DecodeElement(prefab, &token); err != nil {
					return nil, err
				}
				return prefab, nil
			}
			return nil, fmt.Errorf("no prefab loader for %s found", token.Name.Local)
		}
	}
	return nil, fmt.Errorf("could not load prefab %s", resource.Path())
}
