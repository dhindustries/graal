package graal

type Api interface {
	apiMain
	apiHandle
	apiFileSystem
	apiWindow
	apiImage
	apiCamera
	apiRenderer
	apiMesh
	apiTexture
	apiShader
	apiProgram
	apiText
	apiMouse
	apiKeyboard
}

type ApiPrototype struct {
	protoWindow
	protoFileSystem
	protoCamera
	protoRenderer
	protoImage
	protoMesh
	protoTexture
	protoShader
	protoProgram
	protoText
	protoMouse
	protoKeyboard
}

type apiAdapter struct {
	proto     ApiPrototype
	immediate bool
}

var api Api
