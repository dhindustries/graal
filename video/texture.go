package video

import (
	"fmt"

	"github.com/dhindustries/graal"
)

func loadTexture(api *graal.Api, path string) (graal.TextureResource, error) {
	r, err := api.LoadResource(api, "texture", path)
	if err != nil {
		return nil, err
	}
	if v, ok := r.(graal.TextureResource); ok {
		return v, nil
	}
	r.Release()
	return nil, fmt.Errorf("resource %s is not a texture", path)
}
