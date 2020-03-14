package opengl

import (
	"runtime"

	"github.com/dhindustries/graal"
	"github.com/go-gl/gl/v4.3-compatibility/gl"
)

type texture struct {
	graal.Resource
	glid uint32
	img  graal.Image
}

func (texture *texture) Dispose() {
	if !texture.Valid() {
		texture.destroy()
	}
}

func (texture *texture) destroy() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if texture.img != nil {
		texture.img.Release()
		texture.img = nil
	}
	if texture.glid != 0 {
		gl.DeleteTextures(1, &texture.glid)
		texture.glid = 0
	}
}

func textureImageLoader(resource graal.Resource, manager graal.ResourceManager) (graal.Resource, error) {
	builder := builder{}
	img, err := graal.Resources{manager}.LoadImage(resource.Path())
	if err != nil {
		return nil, err
	}
	tex := &texture{resource, 0, img}
	builder.buildTexture(tex)
	return tex, nil
}
