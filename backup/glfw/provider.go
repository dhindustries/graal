package glfw

import (
	"fmt"

	"github.com/dhindustries/graal"
)

type Provider struct{}

func (Provider) Provide(api *graal.Api) error {
	keyboard := newKeyboard()
	mouse := newMouse()
	api.InitSystem = func(api *graal.Api) error {
		api.Invoke(glfwInit)
		return nil
	}
	api.FinitSystem = func(api *graal.Api) {
		api.Invoke(glfwFinit)
	}
	api.InitInput = func(api *graal.Api, wnd graal.Window) error {
		w, ok := wnd.(*Window)
		if !ok {
			return fmt.Errorf("only glfw windows supported")
		}
		w.Window.SetKeyCallback(keyboard.handle)
		w.Window.SetCursorPosCallback(mouse.handleCursor)
		w.Window.SetMouseButtonCallback(mouse.handleButton)
		return nil
	}
	api.UpdateInput = func(api *graal.Api) {
		keyboard.update(api)
		mouse.update(api)
	}
	api.NewWindow = newWindow
	api.IsKeyDown = keyboard.isDown
	api.IsKeyUp = keyboard.isUp
	api.IsKeyPressed = keyboard.isPressed
	api.IsKeyReleased = keyboard.isReleased
	api.IsButtonDown = mouse.isDown
	api.IsButtonUp = mouse.isUp
	api.IsButtonPressed = mouse.isPressed
	api.IsButtonReleased = mouse.isReleased
	api.GetCursorPosition = mouse.cursor

	return nil
}
