package graal

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

type MouseButton = glfw.MouseButton

const (
	MouseButtonLeft  = MouseButton(glfw.MouseButton1)
	MouseButtonRight = MouseButton(glfw.MouseButton2)
)

type Mouse interface {
	IsDown(button MouseButton) bool
	IsUp(button MouseButton) bool
	IsPressed(button MouseButton) bool
	IsReleased(button MouseButton) bool
	Cursor() mgl64.Vec2
}

type apiMouse struct {
	api *Api
}

func (m *apiMouse) IsDown(b MouseButton) bool {
	return m.api.IsButtonDown(m.api, b)
}

func (m *apiMouse) IsUp(b MouseButton) bool {
	return m.api.IsButtonUp(m.api, b)
}

func (m *apiMouse) IsPressed(b MouseButton) bool {
	return m.api.IsButtonPressed(m.api, b)

}

func (m *apiMouse) IsReleased(b MouseButton) bool {
	return m.api.IsButtonReleased(m.api, b)
}

func (m *apiMouse) Cursor() mgl64.Vec2 {
	return m.api.GetCursorPosition(m.api)
}
