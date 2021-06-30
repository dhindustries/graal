package text

import "github.com/dhindustries/graal"

type Text struct {
	graal.Handle
	mesh graal.Mesh
	font Font
	text string
	color graal.Color
}

func (text *Text) Dispose() {
	if text.mesh != nil {
		text.mesh.Release()
		text.mesh = nil
	}
	if text.font != nil {
		text.font.Release()
		text.font = nil
	}
}

func (text *Text) Mesh() graal.Mesh {
	return text.mesh
}

func (text *Text) Color() graal.Color {
	return text.color
}

func (text *Text) SetColor(color graal.Color) {
	text.color = color
}
