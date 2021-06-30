package glfw

import (
	"github.com/dhindustries/graal"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type Library struct {
	window
	mouse
	keyboard
}

func (*Library) Name() string {
	return "glfw v3.2"
}

func (lib *Library) Install(proto *graal.ApiPrototype) error {
	proto.InitWindow = lib.init
	proto.FinishWindow = lib.finish
	proto.UpdateWindow = lib.update
	proto.OpenWindow = lib.window.open
	proto.CloseWindow = lib.window.close
	proto.RenderWindow = lib.window.render
	proto.IsWindowOpen = lib.window.isOpen
	proto.WindowSize = lib.window.size

	proto.MousePosition = lib.mouse.position
	proto.IsKeyDown = lib.keyboard.isKeyDown
	proto.IsKeyUp = lib.keyboard.isKeyUp
	proto.IsKeyReleased = lib.keyboard.isKeyReleased
	proto.IsKeyPressed = lib.keyboard.isKeyPressed
	proto.KeyboardInput = lib.keyboard.input

	return nil
}

func (lib *Library) Init(api graal.Api) error {
	if err := glfw.Init(); err != nil {
		return err
	}
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	return nil
}

func (lib *Library) Finish(api graal.Api) {
	glfw.Terminate()
}

func (lib *Library) init(api graal.Api) error {
	if err := lib.window.init(api); err != nil {
		return err
	}
	lib.keyboard.init(lib.window.window)
	lib.mouse.init(lib.window.window)
	return nil
}

func (lib *Library) finish(api graal.Api) {
	lib.mouse.finish()
	lib.keyboard.finish()
	lib.window.finish(api)
}

func (lib *Library) update(api graal.Api) {
	lib.keyboard.update()
	lib.window.update(api)
}
