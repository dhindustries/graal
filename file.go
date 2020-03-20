package graal

import (
	"io"
)

type File interface {
	Resource
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}
