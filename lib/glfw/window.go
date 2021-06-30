package glfw

import (
	"github.com/dhindustries/graal"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type window struct {
	window *glfw.Window
}

func (wnd *window) init(api graal.Api) error {
	return api.TryInvoke(func(api graal.Api) error {
		var err error
		if wnd.window, err = glfw.CreateWindow(800, 600, "", nil, nil); err != nil {
			return err
		}
		wnd.window.MakeContextCurrent()
		return nil
	})
}

func (wnd *window) finish(api graal.Api) {
	api.Invoke(func(api graal.Api) {
		wnd.window.Destroy()
		wnd.window = nil
	})
}

func (wnd *window) open(api graal.Api, width, height uint, title string) error {
	return api.TryInvoke(func(api graal.Api) error {
		wnd.window.SetSize(int(width), int(height))
		wnd.window.SetTitle(title)
		wnd.window.Show()
		return nil
	})
}

func (wnd *window) close(api graal.Api) {
	api.Invoke(func(api graal.Api) {
		wnd.window.SetShouldClose(true)
		wnd.window.Hide()
	})
}

func (*window) update(api graal.Api) {
	api.Invoke(func(api graal.Api) {
		glfw.PollEvents()
	})
}

func (wnd *window) render(api graal.Api) {
	api.Invoke(func(api graal.Api) {
		wnd.window.SwapBuffers()
	})
}

func (wnd *window) isOpen(api graal.Api) bool {
	return !wnd.window.ShouldClose()
}

func (wnd *window) size(api graal.Api) (width, height uint) {
	w, h := wnd.window.GetSize()
	return uint(w), uint(h)
}
