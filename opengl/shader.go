package opengl

import (
	"fmt"
	"io/ioutil"

	"github.com/dhindustries/graal/memory"

	"github.com/go-gl/gl/v4.6-core/gl"

	"github.com/dhindustries/graal"
)

type shader struct {
	graal.Handle
	id   uint32
	t    graal.ShaderType
	code string
	api  *graal.Api
}

type shaderResource struct {
	graal.Resource
	shader
}

func (sh *shader) Dispose() {
	if v, ok := sh.Handle.(memory.Disposer); ok {
		v.Dispose()
	}
	deleteShader(sh.api, sh.id)
	sh.id = 0
}

func (shader *shader) Type() graal.ShaderType {
	return shader.t
}

func newShader(api *graal.Api, t graal.ShaderType, code string) (graal.Shader, error) {
	sh := &shader{
		Handle: api.NewHandle(api),
		t:      t,
		id:     createShader(api, t),
		code:   code,
		api:    api,
	}
	api.Handle(api, sh)
	compileShader(api, sh.id, sh.code)
	return sh, nil
}

func loadVertexShaderResource(api *graal.Api, r graal.Resource) (graal.Resource, error) {
	return loadShaderResource(api, graal.ShaderTypeVertex, r)
}

func loadFragmentShaderResource(api *graal.Api, r graal.Resource) (graal.Resource, error) {
	return loadShaderResource(api, graal.ShaderTypeFragment, r)
}

func loadShaderResource(api *graal.Api, t graal.ShaderType, r graal.Resource) (graal.Resource, error) {
	f, err := api.LoadFile(api, r.Path())
	if err != nil {
		return nil, err
	}
	defer f.Release()
	code, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	sh := &shaderResource{
		Resource: r,
		shader: shader{
			Handle: r,
			t:      t,
			id:     createShader(api, t),
			code:   gl.GoStr(&code[0]),
			api:    api,
		},
	}
	api.Handle(api, sh)
	compileShader(api, sh.id, sh.code)
	return sh, nil
}

func createShader(api *graal.Api, t graal.ShaderType) uint32 {
	res := make(chan uint32)
	api.Schedule(func() {
		res <- gl.CreateShader(shaderType(t))
	})
	return <-res
}

func deleteShader(api *graal.Api, id uint32) {
	if id != 0 {
		api.Invoke(func() {
			gl.DeleteShader(id)
		})
	}
}

func compileShader(api *graal.Api, id uint32, source string) error {
	if id == 0 {
		return fmt.Errorf("invalid shader to compile")
	}
	cerr := make(chan error)
	code, free := gl.Strs(source + "\x00")
	defer free()
	api.Schedule(func() {
		gl.ShaderSource(id, 1, code, nil)
		gl.CompileShader(id)
		cerr <- getGlError(id, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog, "failed to compile shader")
	})
	return <-cerr
}

func shaderType(t graal.ShaderType) uint32 {
	switch t {
	case graal.ShaderTypeVertex:
		return gl.VERTEX_SHADER
	case graal.ShaderTypeFragment:
		return gl.FRAGMENT_SHADER
	default:
		return 0
	}
}
