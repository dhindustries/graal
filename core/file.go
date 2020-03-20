package core

import (
	"fmt"
	"os"

	"github.com/dhindustries/graal"
	mem "github.com/dhindustries/graal/memory"
)

type file struct {
	graal.Resource
	*os.File
}

func (file *file) Name() string {
	return file.Resource.Name()
}

func (file *file) Dispose() {
	if v, ok := file.Resource.(mem.Disposer); ok {
		v.Dispose()
	}
	file.Close()
}

func loadFileResource(api *graal.Api, r graal.Resource) (graal.Resource, error) {
	f, err := os.Open(r.Path())
	if err != nil {
		return nil, err
	}
	return &file{r, f}, nil
}

func loadFile(api *graal.Api, path string) (graal.File, error) {
	res, err := api.LoadResource(api, "file", path)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(graal.File); ok {
		return v, nil
	}
	return nil, fmt.Errorf("resource %s is not a file", path)
}
