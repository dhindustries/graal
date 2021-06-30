package graal

import (
	"fmt"
	imglib "image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type Image interface {
	Size() (width, height uint)
	Pixel(x, y uint) Color
	SetPixel(x, y uint, color Color)
	Data() []Color
}

type apiImage interface {
	NewImage(width, height uint, data []Color) (Image, error)
	LoadImage(path string) (Image, error)
	Image(image imglib.Image) (Image, error)
}

type protoImage struct {
	NewImage  func(api Api, width, height uint, data []Color) (Image, error)
	LoadImage func(api Api, path string) (Image, error)
	Image     func(api Api, img imglib.Image) (Image, error)
}

func NewImage(width, height uint, data []Color) (Image, error) {
	return api.NewImage(width, height, data)
}

func (api *apiAdapter) NewImage(width, height uint, data []Color) (Image, error) {
	if api.proto.NewImage == nil {
		panic("api.NewImage is not implemented")
	}
	return api.proto.NewImage(api, width, height, data)
}

func newImage(api Api, width, height uint, data []Color) (Image, error) {
	if data != nil && len(data) != int(width*height) {
		return nil, fmt.Errorf("image data is in incorrect size")
	}
	if data == nil {
		data = make([]Color, width*height)
	}
	return &image{width, height, data}, nil
}

func LoadImage(path string) (Image, error) {
	return api.LoadImage(path)
}

func (api *apiAdapter) LoadImage(path string) (Image, error) {
	if api.proto.LoadImage == nil {
		panic("api.LoadImage is not implemented")
	}
	return api.proto.LoadImage(api, path)
}

func loadImage(api Api, path string) (Image, error) {
	file, err := api.ReadFile(path)
	if err != nil {
		return nil, err
	}
	defer api.Release(file)
	img, _, err := imglib.Decode(file)
	if err != nil {
		return nil, err
	}
	return api.Image(img)
}

func (api *apiAdapter) Image(img imglib.Image) (Image, error) {
	if api.proto.Image == nil {
		panic("api.Image is not implemented")
	}
	return api.proto.Image(api, img)
}

func getImage(api Api, img imglib.Image) (Image, error) {
	bounds := img.Bounds()
	width := uint(bounds.Dx())
	height := uint(bounds.Dy())
	data := make([]Color, width*height)
	index := 0
	maxVal := 65535.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			data[index] = Color{
				float64(r) / maxVal,
				float64(g) / maxVal,
				float64(b) / maxVal,
				float64(a) / maxVal,
			}
			index++
		}
	}
	return api.NewImage(width, height, data)
}

type image struct {
	width, height uint
	data          []Color
}

func (img *image) Size() (width, height uint) {
	return img.width, img.height
}

func (img *image) Pixel(x, y uint) Color {
	return img.data[y*img.width+x]
}

func (img *image) SetPixel(x, y uint, color Color) {
	img.data[y*img.width+x] = color
}

func (img *image) Data() []Color {
	return img.data
}
