package video

import (
	"github.com/dhindustries/graal"
)

type Provider struct{}

func (Provider) Provide(api *graal.Api) error {
	api.NewNode = newNode
	api.LoadImage = loadImage
	api.LoadTexture = loadTexture
	api.LoadShader = loadShader
	api.LoadProgram = loadProgram
	api.NewQuad = newQuad
	api.NewOrthoCamera = newOrthoCamera
	api.NewTileset = newTileset
	api.NewTilemap = newTilemap
	if api.SetResourceLoader != nil {
		api.SetResourceLoader(api, "image/*", loadImageResource)
	}
	return nil
}
