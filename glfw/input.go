package glfw

import (
	"fmt"

	"github.com/dhindustries/graal"
)

type Input struct {
	window   *Window
	keyboard *keyboard
}

func (input *Input) Initialize(engine *graal.Engine) error {
	if input.window != nil {
		return fmt.Errorf("input already initialized")
	}
	if window, ok := engine.Window.(*Window); ok {
		input.window = window
		return nil
	}
	return fmt.Errorf("only glfw window supported")
}

func (input *Input) Keyboard() (graal.Keyboard, error) {
	if input.window == nil {
		return nil, fmt.Errorf("input is not initialized")
	}
	if input.keyboard == nil {
		input.keyboard = &keyboard{}
		input.keyboard.currentState = make(map[graal.Key]bool)
		input.keyboard.previousState = make(map[graal.Key]bool)
		input.window.Handle.SetKeyCallback(input.keyboard.handle)
	}
	return input.keyboard, nil
}

func (input *Input) Dispose() {
	if input.keyboard != nil {
		input.window.Handle.SetKeyCallback(nil)
		input.keyboard = nil
	}
	input.window = nil
}

func (input *Input) Update() {
	if input.keyboard != nil {
		input.keyboard.update()
	}
}
