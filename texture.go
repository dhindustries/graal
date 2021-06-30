package graal

type Texture interface {
}

type apiTexture interface {
	NewTexture(image Image) (Texture, error)
	LoadTexture(path string) (Texture, error)
}

type protoTexture struct {
	NewTexture  func(api Api, image Image) (Texture, error)
	LoadTexture func(api Api, path string) (Texture, error)
}

func NewTexture(image Image) (Texture, error) {
	return api.NewTexture(image)
}

func (api *apiAdapter) NewTexture(image Image) (Texture, error) {
	if api.proto.NewTexture == nil {
		panic("api.NewTexture is not implemented")
	}
	return api.proto.NewTexture(api, image)
}

func LoadTexture(path string) (Texture, error) {
	return api.LoadTexture(path)
}

func (api *apiAdapter) LoadTexture(path string) (Texture, error) {
	if api.proto.LoadTexture == nil {
		panic("api.LoadTexture is not implemented")
	}
	return api.proto.LoadTexture(api, path)
}

func loadTexture(api Api, path string) (Texture, error) {
	img, err := api.LoadImage(path)
	if err != nil {
		return nil, err
	}
	api.Release(img)
	return api.NewTexture(img)
}
