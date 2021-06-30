package opengl

import (
	"fmt"

	"github.com/dhindustries/graal"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

const (
	textureName    = "texture"
	worldName      = "model"
	viewName       = "view"
	projectionName = "projection"
	colorName      = "color"
)

type program struct {
	graal.Handle
	id       uint32
	vertex   *vertexShader
	fragment *fragmentShader
}

func newProgram(api graal.Api, vertex graal.VertexShader, fragment graal.FragmentShader) (graal.Program, error) {
	if vertex == nil {
		return nil, fmt.Errorf("vertex shader is not defined")
	}
	if fragment == nil {
		return nil, fmt.Errorf("fragment shader is not defined")
	}
	var vsh *vertexShader
	var fsh *fragmentShader
	var ok bool
	if vsh, ok = vertex.(*vertexShader); !ok {
		return nil, fmt.Errorf("only opengl vertex shader is supported")
	}
	if fsh, ok = fragment.(*fragmentShader); !ok {
		return nil, fmt.Errorf("only opengl fragment shader is supported")
	}
	api.Acquire(vertex)
	api.Acquire(fragment)
	obj := &program{
		id:       glCreateProgram(api),
		vertex:   vsh,
		fragment: fsh,
	}
	glAttachShader(api, obj.id, obj.vertex.id)
	glAttachShader(api, obj.id, obj.fragment.id)
	if err := glLinkProgram(api, obj.id); err != nil {
		api.Release(obj)
		return nil, err
	}
	return obj, nil
}

func (prog *program) Dispose(api graal.Api) {
	if prog.fragment != nil {
		glDetachShader(api, prog.id, prog.fragment.id)
		api.Release(prog.fragment)
		prog.fragment = nil
	}
	if prog.vertex != nil {
		glDetachShader(api, prog.id, prog.vertex.id)
		api.Release(prog.vertex)
		prog.vertex = nil
	}
	if prog.id != 0 {
		glDeleteProgram(api, prog.id)
		prog.id = 0
	}
}

func glCreateProgram(api graal.Api) uint32 {
	var id uint32
	api.Invoke(func(api graal.Api) {
		id = gl.CreateProgram()
	})
	return id
}

func glDeleteProgram(api graal.Api, id uint32) {
	if id != 0 {
		api.Invoke(func(api graal.Api) {
			gl.DeleteProgram(id)
		})
	}
}

func glUseProgram(api graal.Api, id uint32) {
	api.Invoke(func(api graal.Api) {
		gl.UseProgram(id)
	})
}

func glLinkProgram(api graal.Api, id uint32) error {
	return api.TryInvoke(func(api graal.Api) error {
		gl.LinkProgram(id)
		return glError(id, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog, "failed to link program")
	})
}

func glAttachShader(api graal.Api, pid, sid uint32) {
	api.Invoke(func(api graal.Api) {
		gl.AttachShader(pid, sid)
	})
}

func glDetachShader(api graal.Api, pid, sid uint32) {
	api.Invoke(func(api graal.Api) {
		gl.DetachShader(pid, sid)
	})
}

func glSetVec4f(api graal.Api, pid uint32, name string, value mgl32.Vec4) {
	api.Invoke(func(api graal.Api) {
		loc := gl.GetUniformLocation(pid, gl.Str(name+"\x00"))
		gl.Uniform4fv(loc, 1, &value[0])
	})
}

func glSetVec4d(api graal.Api, pid uint32, name string, value mgl64.Vec4) {
	api.Invoke(func(api graal.Api) {
		loc := gl.GetUniformLocation(pid, gl.Str(name+"\x00"))
		gl.Uniform4dv(loc, 1, &value[0])
	})
}

func glSetMat4f(api graal.Api, pid uint32, name string, value mgl32.Mat4) {
	api.Invoke(func(api graal.Api) {
		loc := gl.GetUniformLocation(pid, gl.Str(name+"\x00"))
		gl.UniformMatrix4fv(loc, 1, false, &value[0])
	})
}

func glSetMat4d(api graal.Api, pid uint32, name string, value mgl64.Mat4) {
	api.Invoke(func(api graal.Api) {
		loc := gl.GetUniformLocation(pid, gl.Str(name+"\x00"))
		gl.UniformMatrix4dv(loc, 1, false, &value[0])
	})
}
