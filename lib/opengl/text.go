package opengl

import (
	"sync"
	"sync/atomic"

	"github.com/dhindustries/graal"
	"github.com/go-gl/mathgl/mgl64"
)

func newText(api graal.Api, font graal.Font) (graal.Text, error) {
	api.Acquire(font)
	obj := &text{font: font, chars: map[rune]*textCharacter{}}
	obj.content.Store("test")

	return obj, nil
}

func (renderer *renderer) renderText(api graal.Api, txt *text) {
	txt.mutex.Lock()
	defer txt.mutex.Unlock()
	if content, ok := txt.content.Load().(string); ok {
		x := 0.0
		y := 0.0
		maxWidth := 0.0
		model := renderer.modelTransformation
		prevChar := '\x00'

		if v, ok := txt.maxWidth.Load().(float64); ok {
			maxWidth = v
		}
		for _, char := range content {
			if char == '\n' {
				x = 0
				y += txt.font.LineHeight()
			} else if character := txt.character(api, char); character != nil {
				x -= txt.font.Kerning(prevChar, char)
				if maxWidth > 0 && x+character.advance > maxWidth {
					x = 0
					y += txt.font.LineHeight()
				}
				renderer.applyTransform(api, model.Mul4(mgl64.Translate3D(x, y, 0)))
				renderer.applyTexture(api, character.texture)
				renderer.renderObject(api, character.mesh)
				x += character.advance
				char = prevChar
			}
		}
	}
}

type text struct {
	graal.Handle
	font     graal.Font
	chars    map[rune]*textCharacter
	content  atomic.Value
	mutex    sync.RWMutex
	maxWidth atomic.Value
}

func (t *text) Dispose(api graal.Api) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.font != nil {
		api.Release(t.font)
		t.font = nil
	}
	if t.chars != nil {
		for _, char := range t.chars {
			api.Release(char)
		}
		t.chars = nil
	}
}

func (txt *text) String() string {
	return txt.content.Load().(string)
}

func (txt *text) SetString(content string) {
	txt.mutex.RLock()
	defer txt.mutex.RUnlock()
	txt.content.Store(content)
}

func (txt *text) MaxWidth() float64 {
	return txt.maxWidth.Load().(float64)
}

func (txt *text) SetMaxWidth(value float64) {
	txt.maxWidth.Store(value)
}

func (t *text) character(api graal.Api, char rune) *textCharacter {
	character, ok := t.chars[char]
	if !ok {
		character = newCharacter(api, t.font, char)
		t.chars[char] = character
	}
	return character
}

func newCharacter(api graal.Api, font graal.Font, char rune) *textCharacter {
	glyph := api.FontGlyph(font, char)
	texture, textureLeft, textureTop, textureRight, textureBottom := api.FontTexture(font, char)

	if texture == nil || glyph == nil {
		return nil
	}
	posLeft := glyph.X()
	posTop := glyph.Y() + font.Ascent()
	posRight := posLeft + glyph.Width()
	posBottom := posTop + glyph.Height()

	mesh, _ := api.NewQuad(
		graal.Vertex{
			Position:  mgl64.Vec3{posLeft, posTop, 0},
			TexCoords: mgl64.Vec2{textureLeft, textureTop},
			Color:     mgl64.Vec4{1, 1, 1, 1},
		},
		graal.Vertex{
			Position:  mgl64.Vec3{posRight, posTop, 0},
			TexCoords: mgl64.Vec2{textureRight, textureTop},
			Color:     mgl64.Vec4{1, 1, 1, 1},
		},
		graal.Vertex{
			Position:  mgl64.Vec3{posLeft, posBottom, 0},
			TexCoords: mgl64.Vec2{textureLeft, textureBottom},
			Color:     mgl64.Vec4{1, 1, 1, 1},
		},
		graal.Vertex{
			Position:  mgl64.Vec3{posRight, posBottom, 0},
			TexCoords: mgl64.Vec2{textureRight, textureBottom},
			Color:     mgl64.Vec4{1, 1, 1, 1},
		},
	)
	api.Acquire(texture)
	return &textCharacter{
		mesh:    mesh,
		texture: texture,
		advance: glyph.Advance(),
	}
}

type textCharacter struct {
	texture graal.Texture
	mesh    graal.Mesh
	advance float64
}

func (char *textCharacter) Dispose(api graal.Api) {
	if char.texture != nil {
		api.Release(char.texture)
		char.texture = nil
	}
	if char.mesh != nil {
		api.Release(char.mesh)
		char.mesh = nil
	}
}
