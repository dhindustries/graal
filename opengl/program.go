package opengl

import (
	"fmt"
	"sync"

	"github.com/dhindustries/graal/utils"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"

	"github.com/dhindustries/graal"
	"github.com/dhindustries/graal/memory"
	"github.com/go-gl/gl/v4.6-core/gl"
)

const (
	textureName = "texture"
	worldName   = "world"
	viewName    = "view"
	projName    = "projection"
	colorName   = "color"
)

type program struct {
	graal.Handle
	id   uint32
	vert graal.VertexShader
	frag graal.FragmentShader
	api  *graal.Api
}

type programResource struct {
	graal.Resource
	program
}

func (prog *program) Dispose() {
	if v, ok := prog.Handle.(memory.Disposer); ok {
		v.Dispose()
	}
	prog.SetFragmentShader(nil)
	prog.SetVertexShader(nil)
	deleteProgram(prog.api, prog.id)
	prog.id = 0
}

func (prog *program) SetVec4f(n string, v mgl32.Vec4) {
	setVec4f(prog.api, prog.id, n, v)
}

func (prog *program) SetMat4f(n string, v mgl32.Mat4) {
	setMat4f(prog.api, prog.id, n, v)
}

func (prog *program) SetVertexShader(sh graal.VertexShader) {
	if prog.vert != nil {
		prog.detach(prog.vert)
	}
	if sh != nil {
		prog.attach(sh)
	}
	prog.vert = sh
}

func (prog *program) VertexShader() graal.VertexShader {
	return prog.vert
}

func (prog *program) SetFragmentShader(sh graal.FragmentShader) {
	if prog.frag != nil {
		prog.detach(prog.frag)

	}
	if sh != nil {
		prog.attach(sh)
	}
	prog.frag = sh
}

func (prog *program) FragmentShader() graal.FragmentShader {
	return prog.frag
}

func (prog *program) attach(sh graal.Shader) {
	sh.Acquire()
	if sh, ok := sh.(*shader); ok {
		attachShader(prog.api, prog.id, sh.id)
	}
	if sh, ok := sh.(*shaderResource); ok {
		attachShader(prog.api, prog.id, sh.id)
	}
}

func (prog *program) detach(sh graal.Shader) {
	sh.Release()
	if sh, ok := sh.(*shader); ok {
		detachShader(prog.api, prog.id, sh.id)
	}
	if sh, ok := sh.(*shaderResource); ok {
		detachShader(prog.api, prog.id, sh.id)
	}
}

func (prog *program) Compile() error {
	return linkProgram(prog.api, prog.id)
}

func newProgram(api *graal.Api) (graal.Program, error) {
	prog := &program{
		Handle: api.NewHandle(api),
		id:     createProgram(api),
		api:    api,
	}
	api.Handle(api, prog)
	return prog, nil
}

func loadProgramResource(api *graal.Api, r graal.Resource) (graal.Resource, error) {
	var vert graal.VertexShader
	var frag graal.FragmentShader
	var verr, ferr error
	var wg sync.WaitGroup
	vertName := fmt.Sprintf("%s.vert", r.Path())
	fragName := fmt.Sprintf("%s.frag", r.Path())
	defer func() {
		if vert != nil {
			vert.Release()
		}
		if frag != nil {
			frag.Release()
		}
	}()
	wg.Add(2)
	go func() {
		vert, verr = api.LoadShader(api, graal.ShaderTypeVertex, vertName)
		wg.Done()
	}()
	go func() {
		frag, ferr = api.LoadShader(api, graal.ShaderTypeFragment, fragName)
		wg.Done()
	}()
	wg.Wait()
	if verr != nil {
		return nil, verr
	}
	if ferr != nil {
		return nil, ferr
	}
	prog := &programResource{
		Resource: r,
		program: program{
			Handle: r,
			id:     createProgram(api),
			api:    api,
		},
	}
	prog.SetVertexShader(vert)
	prog.SetFragmentShader(frag)
	if err := prog.Compile(); err != nil {
		prog.Release()
		return nil, err
	}
	return prog, nil
}

func (m *renderer) setModel(api *graal.Api, v mgl64.Mat4) {
	m.setMat4d(api, "model", v)
}

func (m *renderer) setView(api *graal.Api, v mgl64.Mat4) {
	m.setMat4d(api, "view", v)
}

func (m *renderer) setProjection(api *graal.Api, v mgl64.Mat4) {
	m.setMat4d(api, "projection", v)
}

func (m *renderer) setColor(api *graal.Api, v graal.Color) {
	m.setVec4d(api, "color", mgl64.Vec4(v))
}

func (m *renderer) setMat4d(api *graal.Api, name string, v mgl64.Mat4) {
	var id uint32
	switch p := m.prog.(type) {
	case *program:
		id = p.id
	case *programResource:
		id = p.id
	}
	if id != 0 {
		setMat4f(api, id, name, utils.TruncMat4(v))
	}
}

func (m *renderer) setVec4d(api *graal.Api, name string, v mgl64.Vec4) {
	var id uint32
	switch p := m.prog.(type) {
	case *program:
		id = p.id
	case *programResource:
		id = p.id
	}
	if id != 0 {
		setVec4f(api, id, name, utils.TruncVec4(v))
	}
}

func createProgram(api *graal.Api) uint32 {
	res := make(chan uint32)
	defer close(res)
	api.Schedule(func() {
		res <- gl.CreateProgram()
	})
	return <-res
}

func deleteProgram(api *graal.Api, id uint32) {
	if id != 0 {
		api.Invoke(func() {
			gl.DeleteProgram(id)
		})
	}
}

func bindProgram(api *graal.Api, id uint32) {
	api.Invoke(func() {
		gl.UseProgram(id)
	})
}

func linkProgram(api *graal.Api, id uint32) error {
	cerr := make(chan error)
	defer close(cerr)
	api.Schedule(func() {
		gl.LinkProgram(id)
		cerr <- getGlError(id, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog, "failed to link program")
	})
	return <-cerr
}

func attachShader(api *graal.Api, pid, sid uint32) {
	api.Invoke(func() {
		gl.AttachShader(pid, sid)
	})
}

func detachShader(api *graal.Api, pid, sid uint32) {
	api.Invoke(func() {
		gl.DetachShader(pid, sid)
	})
}

func setMat4f(api *graal.Api, pid uint32, name string, val mgl32.Mat4) {
	api.Invoke(func() {
		loc := gl.GetUniformLocation(pid, gl.Str(name+"\x00"))
		gl.UniformMatrix4fv(loc, 1, false, &val[0])
	})
}

func setMat4d(api *graal.Api, pid uint32, name string, val mgl64.Mat4) {
	api.Invoke(func() {
		loc := gl.GetUniformLocation(pid, gl.Str(name+"\x00"))
		gl.UniformMatrix4dv(loc, 1, false, &val[0])
	})
}

func setVec4f(api *graal.Api, pid uint32, name string, val mgl32.Vec4) {
	api.Invoke(func() {
		loc := gl.GetUniformLocation(pid, gl.Str(name+"\x00"))
		gl.Uniform4fv(loc, 1, &val[0])
	})
}

func setVec4d(api *graal.Api, pid uint32, name string, val mgl64.Vec4) {
	api.Invoke(func() {
		loc := gl.GetUniformLocation(pid, gl.Str(name+"\x00"))
		gl.Uniform4dv(loc, 1, &val[0])
	})
}
