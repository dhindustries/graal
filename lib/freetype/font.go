package freetype

import (
	"github.com/dhindustries/graal"
)

type font struct {
	graal.Handle
	glyphs  map[rune]*graal.Glyph
	texture graal.Texture
}

func (f *font) Dispose(api graal.Api) {
	if f.texture != nil {
		api.Release(f.texture)
		f.texture = nil
	}
}

func (lib *Library) loadFont(api graal.Api, path string, size uint) (graal.Font, error) {
	ranges := lib.Ranges
	if ranges == nil {
		ranges = Ranges{{rune(32), rune(127)}}
	}

	img, glyphs, err := loadTrueTypeFont(api, path, size, ranges)
	if err != nil {
		return nil, err
	}
	defer api.Release(img)

	texture, err := api.NewTexture(img)
	if err != nil {
		api.Release(img)
		return nil, err
	}

	return &font{glyphs: glyphs, texture: texture}, nil
}

func (lib *Library) fontGlyph(api graal.Api, f graal.Font, char rune) graal.Glyph {
	if fnt, ok := f.(*font); ok {
		if glyph, ok := fnt.glyphs[char]; ok {
			glyph.Texture = fnt.texture
			return *glyph
		}
	}
	return graal.Glyph{}
}
