package graal

import (
	"io"
	"os"
)

type ReadableFile interface {
	io.Reader
	io.Seeker
	io.Closer
}

type apiFileSystem interface {
	ReadFile(path string) (ReadableFile, error)
}

type protoFileSystem struct {
	ReadFile func(api Api, path string) (ReadableFile, error)
}

func ReadFile(path string) (ReadableFile, error) {
	return api.ReadFile(path)
}

func (api *apiAdapter) ReadFile(path string) (ReadableFile, error) {
	if api.proto.ReadFile == nil {
		panic("api.ReadFile is not implemented")
	}
	return api.proto.ReadFile(api, path)
}

func readFile(api Api, path string) (ReadableFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &readableFile{File: file}, nil
}

type readableFile struct {
	Handle
	*os.File
}

func (file *readableFile) Dispose() {
	file.Close()
}
