package video

import (
	"fmt"

	"github.com/dhindustries/graal"
)

func loadShader(api *graal.Api, t graal.ShaderType, n string) (graal.ShaderResource, error) {
	r, err := api.LoadResource(api, graal.Mime(fmt.Sprintf("shader/%s", t)), n)
	if err != nil {
		return nil, err
	}
	if v, ok := r.(graal.ShaderResource); ok {
		return v, nil
	}
	r.Release()
	return nil, fmt.Errorf("resource %s is not a shader", n)
}
