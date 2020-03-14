package graal

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type Image interface {
	Handle
	Disposable
	Size() (uint, uint)
	At(x, y uint) Color
	Data() []Color
}

const MimeImage = Mime("image/*")

type imageResource struct {
	Resource
	width, height uint
	data          []Color
}

func (image *imageResource) Size() (uint, uint) {
	return image.width, image.height
}

func (image *imageResource) At(x, y uint) Color {
	if x < image.width && y < image.height {
		return image.data[y*image.width+x]
	}
	return Color{}
}

func (image *imageResource) Data() []Color {
	return image.data
}

func (image *imageResource) Dispose() {
	if !image.Valid() {
		image.width = 0
		image.height = 0
		image.data = nil
	}
}

func imageLoader(resource Resource, manager ResourceManager) (Resource, error) {
	file, err := Resources{manager}.LoadFile(resource.Path())
	if err != nil {
		return nil, err
	}
	defer file.Release()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	width := uint(bounds.Dx())
	height := uint(bounds.Dy())
	data := make([]Color, width*height)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		ay := uint(y - bounds.Min.Y)
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			ax := uint(x - bounds.Min.X)
			r, g, b, a := img.At(x, y).RGBA()
			data[ay*width+ax] = Color{float32(r) / 65535.0, float32(g) / 65535.0, float32(b) / 65535.0, float32(a) / 65535.0}
		}
	}
	return &imageResource{resource, width, height, data}, nil
}

func (resources Resources) LoadImage(path string) (Image, error) {
	res, err := resources.Load(MimeImage, path)
	if err != nil {
		return nil, err
	}
	if img, ok := res.(Image); ok {
		return img, nil
	}
	res.Release()
	return nil, fmt.Errorf("resource %s is not a image", path)
}
