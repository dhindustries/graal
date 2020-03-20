package glfw

import (
	"sync"

	"github.com/dhindustries/graal"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type Window struct {
	Window *glfw.Window
	api    *graal.Api
	l      sync.Mutex
}

func newWindow(api *graal.Api) (graal.Window, error) {
	cres := make(chan graal.Window)
	cerr := make(chan error)
	defer close(cres)
	defer close(cerr)
	api.Schedule(func() {
		glfwInit()
		glfw.WindowHint(glfw.Visible, glfw.False)
		glfw.WindowHint(glfw.Resizable, glfw.False)
		glfw.WindowHint(glfw.ContextVersionMajor, 4)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		w, err := glfw.CreateWindow(800, 600, "Graal", nil, nil)
		if err != nil {
			glfwFinit()
			cerr <- err
		} else {
			w.MakeContextCurrent()
			cres <- &Window{Window: w, api: api}
		}
	})
	select {
	case res := <-cres:
		return res, nil
	case err := <-cerr:
		return nil, err
	}
}

func (wnd *Window) Open() error {
	wnd.l.Lock()
	defer wnd.l.Unlock()
	wnd.api.Invoke(func() {
		wnd.Window.Show()
	})
	return nil
}

func (wnd *Window) Close() {
	wnd.l.Lock()
	defer wnd.l.Unlock()
	wnd.api.Invoke(func() {
		wnd.Window.SetShouldClose(true)
		wnd.Window.Hide()
	})
}

func (wnd *Window) Dispose() {
	wnd.l.Lock()
	defer wnd.l.Unlock()
	wnd.api.Invoke(func() {
		wnd.Window.Destroy()
		glfwFinit()
	})
}

func (wnd *Window) IsOpen() bool {
	wnd.l.Lock()
	defer wnd.l.Unlock()
	res := make(chan bool)
	defer close(res)
	wnd.api.Schedule(func() {
		res <- !wnd.Window.ShouldClose()
	})
	return <-res
}

func (wnd *Window) PullMessages() {
	wnd.l.Lock()
	defer wnd.l.Unlock()
	wnd.api.Invoke(func() {
		glfw.PollEvents()
	})
}

func (wnd *Window) Swap() {
	wnd.api.Schedule(func() {
		wnd.Window.SwapBuffers()
	})
}
