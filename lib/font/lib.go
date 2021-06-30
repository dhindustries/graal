package font

import "github.com/dhindustries/graal"

type Library struct {
	Ranges   Ranges
	FontSize uint
	DPI      uint
}

func (*Library) Name() string {
	return "font"
}

func (lib *Library) Install(proto *graal.ApiPrototype) error {
	proto.Font = lib.font
	proto.LoadFont = lib.loadOpenType
	proto.FontGlyph = lib.fontGlyph
	proto.FontTexture = lib.fontTexture

	return nil
}

func (lib *Library) Init(api graal.Api) error {
	if lib.Ranges == nil {
		lib.Ranges = Ranges{{rune(32), rune(128)}}
	}
	if lib.FontSize == 0 {
		lib.FontSize = 36
	}
	if lib.DPI == 0 {
		lib.DPI = 72
	}
	return nil
}

func (*Library) Finish(api graal.Api) {

}
