package graal

type TextureMode string

const (
	TextureModeSharp  = TextureMode("sharp")
	TextureModeSmooth = TextureMode("smooth")
)

type baseTexture interface {
	Draw(img Image) error
	SetMode(mode TextureMode)
	Size() (w, h uint)
}

type Texture interface {
	Handle
	baseTexture
}

type TextureResource interface {
	Resource
	baseTexture
}
