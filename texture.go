package graal

import "fmt"

type Texture interface {
	Handle
	Disposable
}

const MimeTextureImage = Mime("texture/image")

func (resources Resources) LoadTexture(path string) (Texture, error) {
	res, err := resources.Load(MimeTextureImage, path)
	if err != nil {
		return nil, err
	}
	if tex, ok := res.(Texture); ok {
		return tex, nil
	}
	res.Release()
	return nil, fmt.Errorf("Resource %s is not a texture", path)
}
