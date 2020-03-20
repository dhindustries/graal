package opengl

import (
	"sync"

	"github.com/dhindustries/graal"
	"github.com/go-gl/gl/v4.6-core/gl"
)

var glCounter uint = 0
var glLock sync.Mutex

func initOpenGl() error {
	glLock.Lock()
	defer glLock.Unlock()
	if glCounter == 0 {
		return gl.Init()
	}
	glCounter++
	return nil
}

func finitOpenGl() {
	glLock.Lock()
	defer glLock.Unlock()
	if glCounter == 1 {
		gl.Finish()
	}
	glCounter--
}

func initGraphics(api *graal.Api, wnd graal.Window) error {
	cerr := make(chan error)
	defer close(cerr)
	api.Schedule(func() {
		err := initOpenGl()
		if err == nil {
			gl.ActiveTexture(gl.TEXTURE0)
			gl.Enable(gl.TEXTURE_2D)
			gl.Enable(gl.DEPTH_TEST)
			gl.DepthFunc(gl.LESS)
			// gl.Disable(gl.CULL_FACE)
			api.Logf(api, "Initialized OpenGl %s\n", gl.GoStr(gl.GetString(gl.VERSION)))
		}
		cerr <- err
	})
	err := <-cerr
	if err != nil {
		return err
	}
	return nil
}

func finitGraphics(api *graal.Api) {
	api.Invoke(finitOpenGl)
}
