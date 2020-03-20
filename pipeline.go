package graal

type Pipeline interface {
	Disposable
	SetCamera(camera Camera)
	Camera() Camera
	SetProgram(program Program)
	Program() Program
}
