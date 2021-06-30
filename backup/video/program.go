package video

import (
	"fmt"

	"github.com/dhindustries/graal"
)

func loadProgram(api *graal.Api, n string) (graal.ProgramResource, error) {
	r, err := api.LoadResource(api, "program", n)
	if err != nil {
		return nil, err
	}
	if v, ok := r.(graal.ProgramResource); ok {
		return v, nil
	}
	r.Release()
	return nil, fmt.Errorf("resource %s is not a program", n)
}
