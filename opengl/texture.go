package opengl

import (
	"github.com/dhindustries/graal"
	"github.com/dhindustries/graal/memory"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type texture struct {
	graal.Handle
	id   uint32
	w, h uint
	api  *graal.Api
}

type textureResource struct {
	graal.Resource
	texture
}

func (t *texture) Size() (uint, uint) {
	return t.w, t.h
}

func (t *texture) Draw(img graal.Image) error {
	w, h := img.Size()
	fillTexture(t.api, t.id, w, h, img.Data())
	return nil
}

func (t *texture) SetMode(m graal.TextureMode) {
	switch m {
	case graal.TextureModeSharp:
		setTextureMode(t.api, t.id, gl.NEAREST)
	case graal.TextureModeSmooth:
		setTextureMode(t.api, t.id, gl.LINEAR)
	}
}

func (t *texture) Dispose() {
	memory.Dispose(t.Handle)
	deleteTexture(t.api, t.id)
	t.id = 0
}

func newTexture(api *graal.Api, w, h uint) (graal.Texture, error) {
	tex := &texture{
		Handle: api.NewHandle(api),
		id:     createTexture(api),
		w:      w, h: h,
		api: api,
	}
	api.Handle(api, tex)
	return tex, nil
}

func loadTextureResource(api *graal.Api, r graal.Resource) (graal.Resource, error) {
	img, err := api.LoadImage(api, r.Path())
	if err != nil {
		return nil, err
	}
	defer img.Release()
	w, h := img.Size()
	tex := &textureResource{
		Resource: r,
		texture: texture{
			Handle: r,
			id:     createTexture(api),
			w:      w, h: h,
			api: api,
		},
	}
	fillTexture(api, tex.id, w, h, img.Data())
	return tex, nil
}

func createTexture(api *graal.Api) uint32 {
	res := make(chan uint32)
	api.Schedule(func() {
		var id uint32
		gl.GenTextures(1, &id)
		gl.BindTexture(gl.TEXTURE_2D, id)
		defer gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
		res <- id
	})
	return <-res
}

func bindTexture(api *graal.Api, tex graal.Texture) {
	api.Invoke(func() {
		var id uint32
		switch v := tex.(type) {
		case *texture:
			id = v.id
		case *textureResource:
			id = v.id
		}
		gl.BindTexture(gl.TEXTURE_2D, id)
	})
}

func deleteTexture(api *graal.Api, id uint32) {
	if id != 0 {
		api.Schedule(func() {
			gl.DeleteTextures(1, &id)
		})
	}
}

func setTextureMode(api *graal.Api, id uint32, m int32) {
	if id != 0 {
		api.Invoke(func() {
			gl.BindTexture(gl.TEXTURE_2D, id)
			defer gl.BindTexture(gl.TEXTURE_2D, 0)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, m)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, m)
		})
	}
}

func fillTexture(api *graal.Api, id uint32, w, h uint, data []graal.Color) {
	api.Invoke(func() {
		gl.BindTexture(gl.TEXTURE_2D, id)
		defer gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.FLOAT, gl.Ptr(data))
		gl.GenerateMipmap(gl.TEXTURE_2D)
	})
}

func (m *renderer) pushTexture(t graal.Texture) {
	var id uint32
	switch v := t.(type) {
	case *texture:
		id = v.id
	case *textureResource:
		id = v.id
	}
	gl.BindTexture(gl.TEXTURE_2D, id)
	// if m.textureStack != nil {
	// 	m.textureStack.push(id)
	// }
}

func (m *renderer) popTexture(n int) {
	// if m.textureStack != nil {
	// 	gl.BindTexture(gl.TEXTURE_2D, m.textureStack.popn(n))
	// } else if n > 0 {
	gl.BindTexture(gl.TEXTURE_2D, 0)
	// }
}

func (m *renderer) rebindTexture() {
	// gl.BindTexture(gl.TEXTURE_2D, m.textureStack.top())
}
