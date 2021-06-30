package opengl

import (
	"fmt"
	"sync"

	"github.com/dhindustries/graal"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type shader struct {
	graal.Handle
	id    uint32
	mutex sync.Mutex
}

type fragmentShader struct {
	*shader
}

type vertexShader struct {
	*shader
}

func newVertexShader(api graal.Api, source string) (graal.VertexShader, error) {
	sh, err := newShader(api, gl.VERTEX_SHADER, source)
	if err != nil {
		return nil, err
	}
	return &vertexShader{sh}, nil
}

func newFragmentShader(api graal.Api, source string) (graal.FragmentShader, error) {
	sh, err := newShader(api, gl.FRAGMENT_SHADER, source)
	if err != nil {
		return nil, err
	}
	return &fragmentShader{sh}, nil
}

func newShader(api graal.Api, typ uint32, source string) (*shader, error) {
	obj := &shader{
		id: glCreateShader(api, typ),
	}
	if err := glCompileShader(api, obj.id, source); err != nil {
		api.Release(obj)
		return nil, err
	}
	return obj, nil
}

func (sh *shader) Dispose(api graal.Api) {
	sh.mutex.Lock()
	id := sh.id
	sh.id = 0
	sh.mutex.Unlock()
	go glDeleteShader(api, id)
}

func glCreateShader(api graal.Api, typ uint32) uint32 {
	var id uint32
	api.Invoke(func(api graal.Api) {
		id = gl.CreateShader(typ)
	})
	return id
}

func glDeleteShader(api graal.Api, id uint32) {
	if id != 0 {
		api.Invoke(func(api graal.Api) {
			gl.DeleteShader(id)
		})
	}
}

func glCompileShader(api graal.Api, id uint32, source string) error {
	if id == 0 {
		return fmt.Errorf("invalid shader")
	}
	code, free := gl.Strs(source + "\x00")
	defer free()
	return api.TryInvoke(func(api graal.Api) error {
		gl.ShaderSource(id, 1, code, nil)
		gl.CompileShader(id)
		return glError(id, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog, "failed to compile shader")
	})
}
