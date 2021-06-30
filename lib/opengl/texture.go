package opengl

import (
	"sync"

	"github.com/dhindustries/graal"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type texture struct {
	id    uint32
	mutex sync.Mutex
}

func newTexture(api graal.Api, image graal.Image) (graal.Texture, error) {
	tex := &texture{
		id: glCreateTexture(api),
	}
	width, height := image.Size()
	data := colorTrunc(image.Data())
	glUpdateTexture(api, tex.id, int32(width), int32(height), data)
	return tex, nil
}

func (tex *texture) Dispose(api graal.Api) {
	tex.mutex.Lock()
	id := tex.id
	tex.id = 0
	tex.mutex.Unlock()
	glDeleteTexture(api, id)
}

func glCreateTexture(api graal.Api) uint32 {
	var id uint32
	api.Invoke(func(api graal.Api) {
		gl.GenTextures(1, &id)
		gl.BindTexture(gl.TEXTURE_2D, id)
		defer gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	})
	return id
}

func glDeleteTexture(api graal.Api, id uint32) {
	if id != 0 {
		api.Invoke(func(api graal.Api) {
			gl.DeleteTextures(1, &id)
		})
	}
}

func glBindTexture(api graal.Api, id uint32) {
	api.Invoke(func(api graal.Api) {
		gl.BindTexture(gl.TEXTURE_2D, id)
	})
}

func glUpdateTexture(api graal.Api, id uint32, width, height int32, data []mgl32.Vec4) {
	api.Invoke(func(api graal.Api) {
		gl.BindTexture(gl.TEXTURE_2D, id)
		defer gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.FLOAT, gl.Ptr(data))
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.GenerateMipmap(gl.TEXTURE_2D)
	})
}
