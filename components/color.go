package components

import "github.com/dhindustries/graal"

type Colored struct {
	color graal.Color
}

func (component *Colored) Color() graal.Color {
	return component.color
}

func (component *Colored) SetColor(color graal.Color) {
	component.color = color
}
