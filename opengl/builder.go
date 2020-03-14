package opengl

import (
	"runtime"

	"github.com/go-gl/gl/v4.3-compatibility/gl"
)

type builder struct{}

func (*builder) buildTexture(tex *texture) {
	if tex.glid == 0 && tex.img != nil {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
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
		runtime.SetFinalizer(tex, func(tex *texture) {
			tex.destroy()
		})
	}
}
