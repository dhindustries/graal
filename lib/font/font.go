package font

import (
	"image"
	"math"

	"github.com/dhindustries/graal"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func (lib *Library) font(api graal.Api, face font.Face) (graal.Font, error) {
	metrics := face.Metrics()
	runes := lib.Ranges.Runes()
	maxGlyphWidth := uint32(0)
	maxGlyphHeight := uint32(metrics.Height >> 6)
	for _, ch := range runes {
		bounds, _, ok := face.GlyphBounds(ch)
		if ok {
			width := uint32((bounds.Max.X - bounds.Min.X) >> 6)
			maxGlyphWidth = max(maxGlyphWidth, width+uint32(lib.FontSize/2))
		}
	}

	glyphCount := len(runes)
	glyphCols := int(math.Sqrt(float64(glyphCount))) + 1
	glyphRows := glyphCount/glyphCols + 1
	imageWidth := pow2(uint32(glyphCols) * maxGlyphWidth)
	imageHeight := pow2(uint32(glyphRows) * maxGlyphHeight)

	index := 0
	target := image.NewRGBA(image.Rect(0, 0, int(imageWidth), int(imageHeight)))
	drawer := font.Drawer{
		Face: face,
		Src:  image.White,
		Dst:  target,
	}

	x := 0
	y := int(metrics.Ascent >> 6)
	size := float64(lib.FontSize)
	textureWidth := float64(imageWidth)
	textureHeight := float64(imageHeight)
	glyphs := make(map[rune]*fontGlyph, glyphCount)
	for _, ch := range runes {
		start := fixed.P(x, y)
		drawer.Dot = start
		drawer.DrawString(string(ch))

		if bounds, advance, ok := face.GlyphBounds(ch); ok {
			glyph := &fontGlyph{}
			glyph.left = float64((bounds.Min.X+start.X)>>6) / textureWidth
			glyph.top = float64((bounds.Min.Y+start.Y)>>6) / textureHeight
			glyph.right = float64((bounds.Max.X+start.X)>>6) / textureWidth
			glyph.bottom = float64((bounds.Max.Y+start.Y)>>6) / textureHeight
			glyph.x = float64(bounds.Min.X>>6) / size
			glyph.y = float64(bounds.Min.Y>>6) / size
			glyph.width = float64((bounds.Max.X-bounds.Min.X)>>6) / size
			glyph.height = float64((bounds.Max.Y-bounds.Min.Y)>>6) / size
			glyph.advance = float64(advance>>6) / size
			glyphs[ch] = glyph
		}

		index++
		if index%glyphCols == 0 {
			x = 0
			y += int(maxGlyphHeight)
		} else {
			x += int(maxGlyphWidth)
		}
	}

	img, err := api.Image(target)
	if err != nil {
		face.Close()
		return nil, err
	}
	defer api.Release(img)

	texture, err := api.NewTexture(img)
	if err != nil {
		face.Close()
		return nil, err
	}

	return &fontFace{
		face:       face,
		glyphs:     glyphs,
		texture:    texture,
		maxWidth:   float64(maxGlyphWidth) / size,
		maxHeight:  float64(maxGlyphHeight) / size,
		lineHeight: float64(metrics.Height>>6) / size,
		ascent:     float64(metrics.Ascent>>6) / size,
		descent:    float64(metrics.Descent>>6) / size,
		size:       size,
	}, nil
}

func (*Library) fontGlyph(api graal.Api, font graal.Font, char rune) graal.Glyph {
	if f, ok := font.(*fontFace); ok {
		if g, ok := f.glyphs[char]; ok {
			return g
		}
	}
	return nil
}

func (*Library) fontTexture(api graal.Api, font graal.Font, char rune) (texture graal.Texture, left, top, right, bottom float64) {
	if f, ok := font.(*fontFace); ok {
		if g, ok := f.glyphs[char]; ok {
			return f.texture, g.left, g.top, g.right, g.bottom
		}
	}
	return nil, 0, 0, 0, 0
}

type fontFace struct {
	face       font.Face
	glyphs     map[rune]*fontGlyph
	texture    graal.Texture
	maxWidth   float64
	maxHeight  float64
	lineHeight float64
	ascent     float64
	descent    float64
	size       float64
}

func (font *fontFace) Dispise(api graal.Api) {
	if font.face != nil {
		font.face.Close()
		font.face = nil
	}
	if font.texture != nil {
		api.Release(font.texture)
		font.texture = nil
	}
	if font.glyphs != nil {
		font.glyphs = nil
	}
}

func (font *fontFace) MaxWidth() float64 {
	return font.maxWidth
}

func (font *fontFace) MaxHeight() float64 {
	return font.maxHeight
}

func (font *fontFace) LineHeight() float64 {
	return font.lineHeight
}

func (font *fontFace) Ascent() float64 {
	return font.ascent
}

func (font *fontFace) Descent() float64 {
	return font.descent
}

func (font *fontFace) Kerning(a, b rune) float64 {
	return float64(font.face.Kern(a, b)>>6) / font.size
}

type fontGlyph struct {
	left    float64
	top     float64
	right   float64
	bottom  float64
	x       float64
	y       float64
	width   float64
	height  float64
	advance float64
}

func (glyph *fontGlyph) X() float64 {
	return glyph.x
}

func (glyph *fontGlyph) Y() float64 {
	return glyph.y
}

func (glyph *fontGlyph) Width() float64 {
	return glyph.width
}

func (glyph *fontGlyph) Height() float64 {
	return glyph.height
}

func (glyph *fontGlyph) Advance() float64 {
	return glyph.advance
}
