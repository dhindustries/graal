package freetype

import "github.com/dhindustries/graal"

type Library struct {
	Ranges Ranges
}

func (*Library) Name() string {
	return "freetype"
}

func (lib *Library) Install(proto *graal.ApiPrototype) error {
	proto.LoadFont = lib.loadFont
	proto.FontGlyph = lib.fontGlyph

	return nil
}

func (*Library) Init(api graal.Api) error {
	return nil
}

func (*Library) Finish(api graal.Api) {

}
