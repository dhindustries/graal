package graal

type Program interface {
	Disposable
	SetWorld(mat Mat4x4)
	SetProjection(mat Mat4x4)
	SetView(mat Mat4x4)
	SetTexture(tex Texture)
}
