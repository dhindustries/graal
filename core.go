package graal

type core struct{}

func (*core) Name() string {
	return "core"
}

func (*core) Install(proto *ApiPrototype) error {
	proto.ReadFile = readFile
	proto.Image = getImage
	proto.NewImage = newImage
	proto.LoadImage = loadImage
	proto.LoadTexture = loadTexture
	proto.NewQuad = newQuad
	proto.NewOrthoCamera = newOrthoCamera
	proto.NewPerspectiveCamera = newPerspectiveCamera

	return nil
}

func (*core) Init(api Api) error {
	return nil
}

func (*core) Finish(api Api) {

}
