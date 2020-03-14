package opengl

import (
	"runtime"

	"github.com/dhindustries/graal"
	"github.com/go-gl/gl/v4.3-compatibility/gl"
)

type texture struct {
	graal.Resource
	glid  uint32
	img   graal.Image
	queue *graal.Queue
}

func (texture *texture) Dispose() {
	if !texture.Valid() {
		texture.destroy()
	}
}

func (texture *texture) destroy() {
	if texture.img != nil {
		texture.img.Release()
		texture.img = nil
	}
	if texture.glid != 0 {
		// texture.queue.Exec(func() {
		runtime.LockOSThread()
		gl.DeleteTextures(1, &texture.glid)
		texture.glid = 0
		runtime.UnlockOSThread()
		// })
	}
}

func (builder *builder) buildTexture(tex *texture) {
	if tex.glid == 0 && tex.img != nil {
		// builder.queue.Exec(func() {
		runtime.LockOSThread()
		gl.GenTextures(1, &tex.glid)
		gl.BindTexture(gl.TEXTURE_2D, tex.glid)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		w, h := tex.img.Size()
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.FLOAT, gl.Ptr(tex.img.Data()))
		gl.GenerateMipmap(gl.TEXTURE_2D)
		tex.img.Release()
		tex.img = nil
		runtime.UnlockOSThread()
		// })
		runtime.SetFinalizer(tex, func(tex *texture) {
			tex.destroy()
		})
	}
}

func (loader *loader) loadTexture(resource graal.Resource, manager graal.ResourceManager) (graal.Resource, error) {
	img, err := graal.Resources{manager}.LoadImage(resource.Path())
	if err != nil {
		return nil, err
	}
	tex := &texture{Resource: resource, glid: 0, img: img, queue: loader.queue}
	loader.builder.buildTexture(tex)
	return tex, nil
}
