package graal

type Material interface {
	Handle
	Disposable
	Texture() Texture
}
