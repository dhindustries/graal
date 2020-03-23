package graal

import (
	"github.com/go-gl/mathgl/mgl64"
)

type MemoryApi struct {
	NewHandle func(*Api) Handle
	Handle    func(*Api, Handle)
	Cleanup   func(*Api)
}

type RuntimeApi struct {
	Schedule func(func())
	Invoke   func(func())
}

type LogApi struct {
	Logf func(*Api, string, ...interface{})
}

type ConfigApi struct {
	HasConfigValue  func(*Api, string) bool
	SetConfigValue  func(*Api, string, interface{})
	GetConfigValue  func(*Api, string, interface{}) (interface{}, bool)
	GetConfigInt    func(*Api, string, int) (int, bool)
	GetConfigUint   func(*Api, string, uint) (uint, bool)
	GetConfigString func(*Api, string, string) (string, bool)
	GetConfigFloat  func(*Api, string, float64) (float64, bool)
	GetConfigBool   func(*Api, string, bool) (bool, bool)
}

type ResourceApi struct {
	NewResource       func(*Api, Mime, string) (Resource, error)
	LoadResource      func(*Api, Mime, string) (Resource, error)
	GetRelativePath   func(*Api, Resource, string) string
	SetResourceLoader func(*Api, Mime, func(*Api, Resource) (Resource, error))
}

type PrefabApi struct {
	LoadPrefab      func(*Api, string) (Prefab, error)
	SetPrefabLoader func(*Api, string, string, PrefabLoader)
}

type WindowApi struct {
	NewWindow  func(*Api) (Window, error)
	SwapWindow func(*Api, Window)
}

type FileApi struct {
	LoadFile func(*Api, string) (File, error)
}

type KeyboardApi struct {
	IsKeyUp       func(*Api, Key) bool
	IsKeyDown     func(*Api, Key) bool
	IsKeyPressed  func(*Api, Key) bool
	IsKeyReleased func(*Api, Key) bool
}

type MouseApi struct {
	IsButtonUp        func(*Api, MouseButton) bool
	IsButtonDown      func(*Api, MouseButton) bool
	IsButtonPressed   func(*Api, MouseButton) bool
	IsButtonReleased  func(*Api, MouseButton) bool
	GetCursorPosition func(*Api) mgl64.Vec2
}

type InputApi struct {
	KeyboardApi
	MouseApi
	UpdateInput func(*Api)
	InitInput   func(*Api, Window) error
	FinitInput  func(*Api)
}

type RenderApi struct {
	RenderEnqueue func(*Api, interface{}, mgl64.Mat4)
	RenderCommit  func(*Api)
	SetCamera     func(*Api, Camera)
}

type PolygonApi struct {
	NewQuad  func(api *Api, left, right, top, bottom float64) (Shape, error)
	NewMesh  func(*Api, []Vertex) (Mesh, error)
	LoadMesh func(*Api, string) (Mesh, error)
}

type SceneApi struct {
	NewNode   func(*Api) (Node, error)
	SetParent func(*Api, interface{}, Node)
	GetParent func(*Api, interface{}) Node

	NewOrthoCamera func(*Api) (OrthoCamera, error)
	NewTileset     func(*Api) (Tileset, error)
	NewTilemap     func(*Api) (Tilemap, error)
}

type ImageApi struct {
	NewImage    func(*Api, uint, uint) (Image, error)
	LoadImage   func(*Api, string) (ImageResource, error)
	NewTexture  func(*Api, uint, uint) (Texture, error)
	LoadTexture func(*Api, string) (TextureResource, error)
	FillTexture func(*Api, Texture, Image) error
}

type ShaderApi struct {
	NewProgram       func(*Api) (Program, error)
	NewShader        func(*Api, ShaderType, string) (Shader, error)
	AttachShader     func(*Api, Program, Shader) error
	DetachShader     func(*Api, Program, Shader) error
	LoadShader       func(*Api, ShaderType, string) (ShaderResource, error)
	LoadProgram      func(*Api, string) (ProgramResource, error)
	SetShaderMatrix  func(*Api, Program, string, mgl64.Mat4)
	SetShaderTexture func(*Api, Program, string, Texture)
	SetShaderVec4    func(*Api, Program, string, mgl64.Vec4)
	BindProgram      func(*Api, Program)
}

type SystemApi struct {
	WindowApi
	InputApi
	ResourceApi
	PrefabApi
	FileApi
	InitSystem  func(*Api) error
	FinitSystem func(*Api)
}

type GraphicsApi struct {
	RenderApi
	PolygonApi
	ImageApi
	ShaderApi
	SceneApi
	InitGraphics  func(*Api, Window) error
	FinitGraphics func(*Api)
}

type Api struct {
	MemoryApi
	RuntimeApi
	ConfigApi
	LogApi
	SystemApi
	GraphicsApi
}

type ApiProvider interface {
	Provide(*Api) error
}

func Fork(api *Api) *Api {
	f := *api
	return &f
}
