package prefab

import (
	"fmt"

	"github.com/dhindustries/graal"
)

type basePrefab struct {
	graal.Resource
}

func (*basePrefab) Spawn() (graal.Handle, error) {
	panic("not implemented")
}

func loadPrefab(api *graal.Api, path string) (graal.Prefab, error) {
	res, err := api.LoadResource(api, "prefab", path)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(graal.Prefab); ok {
		return v, nil
	}
	res.Release()
	return nil, fmt.Errorf("resource %s is not a prefab", path)
}
