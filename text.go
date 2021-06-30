package graal

import "golang.org/x/image/font"

type Glyph interface {
	X() float64
	Y() float64
	Width() float64
	Height() float64
	Advance() float64
}

type Font interface {
	MaxWidth() float64
	MaxHeight() float64
	LineHeight() float64
	Ascent() float64
	Descent() float64
	Kerning(a, b rune) float64
}

type Text interface {
	String() string
	SetString(content string)
	MaxWidth() float64
	SetMaxWidth(value float64)
}

type apiText interface {
	Font(face font.Face) (Font, error)
	LoadFont(path string) (Font, error)
	FontGlyph(font Font, char rune) Glyph
	FontTexture(font Font, char rune) (texture Texture, left, top, right, bottom float64)
	NewText(font Font) (Text, error)
}

type protoText struct {
	Font        func(api Api, face font.Face) (Font, error)
	LoadFont    func(api Api, path string) (Font, error)
	FontGlyph   func(api Api, font Font, char rune) Glyph
	FontTexture func(api Api, font Font, char rune) (texture Texture, left, top, right, bottom float64)
	NewText     func(api Api, font Font) (Text, error)
}

func (api *apiAdapter) Font(face font.Face) (Font, error) {
	if api.proto.Font == nil {
		panic("api.Font is not implemented")
	}
	return api.proto.Font(api, face)
}

func LoadFont(path string) (Font, error) {
	return api.LoadFont(path)
}

func (api *apiAdapter) LoadFont(path string) (Font, error) {
	if api.proto.LoadFont == nil {
		panic("api.LoadFont is not implemented")
	}
	return api.proto.LoadFont(api, path)
}

func NewText(font Font) (Text, error) {
	return api.NewText(font)
}

func (api *apiAdapter) NewText(font Font) (Text, error) {
	if api.proto.NewText == nil {
		panic("api.NewText is not implemented")
	}
	return api.proto.NewText(api, font)
}

func (api *apiAdapter) FontGlyph(font Font, char rune) Glyph {
	if api.proto.FontGlyph == nil {
		panic("api.FontGlyph is not implemented")
	}
	return api.proto.FontGlyph(api, font, char)
}

func (api *apiAdapter) FontTexture(font Font, char rune) (texture Texture, left, top, right, bottom float64) {
	if api.proto.FontTexture == nil {
		panic("api.FontTexture is not implemented")
	}
	return api.proto.FontTexture(api, font, char)
}
