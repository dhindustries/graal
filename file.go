package graal

import (
	"fmt"
	"io"
	"os"
)

type File interface {
	Resource
	Disposable
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}

const MimeFile = Mime("file/*")

type fileResource struct {
	Resource
	*os.File
}

func (file *fileResource) Name() string {
	return file.Resource.Name()
}

func (file *fileResource) Dispose() {
	file.Close()
}

func fileLoader(resource Resource, manager ResourceManager) (Resource, error) {
	file, err := os.Open(resource.Path())
	if err != nil {
		return nil, err
	}

	return &fileResource{resource, file}, nil
}

func (resources Resources) LoadFile(name string) (File, error) {
	res, err := resources.Load(MimeFile, name)
	if err != nil {
		return nil, err
	}
	if file, ok := res.(File); ok {
		return file, nil
	}
	res.Release()
	return nil, fmt.Errorf("resource %s is not a file", name)
}
