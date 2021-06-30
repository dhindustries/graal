package video

import (
	"fmt"
	"sync"

	imglib "image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/dhindustries/graal"
	mem "github.com/dhindustries/graal/memory"
)

type image struct {
	graal.Handle
	w, h uint
	d    []graal.Color
	l    sync.RWMutex
}

type imageResource struct {
	image
	graal.Resource
}

func (img *image) Size() (uint, uint) {
	return img.w, img.h
}

func (img *image) Get(x, y uint) graal.Color {
	if x < img.w && y < img.h {
		img.l.RLock()
		defer img.l.RUnlock()
		return img.d[y*img.w+x]
	}
	return graal.Color{}
}

func (img *image) Set(x, y uint, c graal.Color) {
	if x < img.w && y < img.h {
		img.l.Lock()
		defer img.l.Unlock()
		img.d[y*img.w+x] = c
	}
}

func (img *image) Data() []graal.Color {
	img.l.RLock()
	defer img.l.RUnlock()
	return img.d
}

func (img *image) Dispose() {
	img.l.Lock()
	defer img.l.Unlock()
	if v, ok := img.Handle.(mem.Disposer); ok {
		v.Dispose()
	}
	img.w = 0
	img.h = 0
	img.d = nil
}

func newImage(api *graal.Api, w, h uint) (graal.Image, error) {
	r := &image{
		Handle: api.NewHandle(api),
		w:      w, h: h,
		d: make([]graal.Color, w*h),
	}
	api.Handle(api, r)
	return r, nil
}

func loadImageResource(api *graal.Api, r graal.Resource) (graal.Resource, error) {
	f, err := api.LoadFile(api, r.Path())
	if err != nil {
		return nil, err
	}
	defer f.Release()
	img, _, err := imglib.Decode(f)
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	w := uint(bounds.Dx())
	h := uint(bounds.Dy())
	data := make([]graal.Color, w*h)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		ay := uint(y - bounds.Min.Y)
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			ax := uint(x - bounds.Min.X)
			r, g, b, a := img.At(x, y).RGBA()
			data[ay*w+ax] = graal.Color{
				float64(r) / 65535.0,
				float64(g) / 65535.0,
				float64(b) / 65535.0,
				float64(a) / 65535.0,
			}
		}
	}
	return &imageResource{
		Resource: r,
		image: image{
			Handle: r,
			w:      w, h: h,
			d: data,
		},
	}, nil
}

func loadImage(api *graal.Api, path string) (graal.ImageResource, error) {
	r, err := api.LoadResource(api, "image", path)
	if err != nil {
		return nil, err
	}
	if v, ok := r.(graal.ImageResource); ok {
		return v, nil
	}
	r.Release()
	return nil, fmt.Errorf("resource %s is not image resource", path)
}
