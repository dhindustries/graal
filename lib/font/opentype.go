package font

import (
	"io"
	"io/ioutil"

	"github.com/dhindustries/graal"
	"golang.org/x/image/font/opentype"
)

func (lib *Library) loadOpenType(api graal.Api, path string) (graal.Font, error) {
	file, err := api.ReadFile(path)
	if err != nil {
		return nil, err
	}
	defer api.Release(file)
	font, err := lib.parseOpenType(file)
	if err != nil {
		return nil, err
	}
	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: float64(lib.FontSize),
		DPI:  float64(lib.DPI),
	})
	if err != nil {
		return nil, err
	}
	return api.Font(face)
}

func (*Library) parseOpenType(file graal.ReadableFile) (*opentype.Font, error) {
	if f, ok := file.(io.ReaderAt); ok {
		return opentype.ParseReaderAt(f)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return opentype.Parse(data)
}
