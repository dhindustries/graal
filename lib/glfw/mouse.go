package glfw

import (
	"github.com/dhindustries/graal"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type mouse struct {
	window *glfw.Window
}

func (mouse *mouse) init(window *glfw.Window) {
	mouse.window = window
}

func (mouse *mouse) finish() {
	mouse.window = nil
}

func (mouse *mouse) position(api graal.Api) (uint, uint) {
	x, y := mouse.window.GetCursorPos()
	return uint(x), uint(y)
}
