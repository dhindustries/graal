package opengl

import (
	"sync"
	"unsafe"

	"github.com/dhindustries/graal"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func newMesh(api graal.Api, vxs []graal.Vertex, idxs []uint) (graal.Mesh, error) {
	obj := &mesh{
		vao: glCreateVertexArray(api),
		vba: glCreateBufferArray(api),
	}
	if vxs != nil {
		var vc int
		if idxs == nil {
			vc = len(vxs)
		} else {
			vc = len(idxs)
		}
		obj.from = 0
		obj.to = int32(vc)
		glUpdateVertexBuffers(api, obj.vao, obj.vba, vxs, idxs)
	}
	return obj, nil
}

const (
	positionLocation  = 0
	normalLocation    = 1
	texcoordsLocation = 2
	colorLocation     = 3
)

type mesh struct {
	graal.Handle
	vao      uint32
	vba      [4]uint32
	from, to int32
	mutex    sync.Mutex
}

type meshAttribute struct {
	location uint32
	size     uint32
	dim      int32
	typ      uint32
	norm     bool
	data     unsafe.Pointer
}

func (mesh *mesh) Dispose(api graal.Api) {
	mesh.mutex.Lock()
	vao := mesh.vao
	vba := mesh.vba
	mesh.vao = 0
	mesh.vba = [4]uint32{}
	mesh.from = 0
	mesh.to = 0
	mesh.mutex.Unlock()
	glDeleteBufferArray(api, vba)
	glDeleteVertexArray(api, vao)
}

func (renderer *renderer) renderMesh(api graal.Api, obj *mesh) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	glDrawVertexArray(api, obj.vao, gl.TRIANGLES, obj.from, obj.to)
}

func glCreateVertexArray(api graal.Api) uint32 {
	var vao uint32
	api.Invoke(func(api graal.Api) {
		gl.CreateVertexArrays(1, &vao)
	})
	return vao
}

func glCreateBufferArray(api graal.Api) [4]uint32 {
	vba := [4]uint32{}
	api.Invoke(func(api graal.Api) {
		gl.CreateBuffers(4, &vba[0])
	})
	return vba
}

func glUpdateVertexBuffers(api graal.Api, vao uint32, vba [4]uint32, vxs []graal.Vertex, idxs []uint) {
	if vxs == nil {
		return
	}
	var vertexCount int
	if idxs == nil {
		vertexCount = len(vxs)
	} else {
		vertexCount = len(idxs)
	}
	positions := make([]mgl32.Vec3, vertexCount)
	normals := make([]mgl32.Vec3, vertexCount)
	texcoords := make([]mgl32.Vec2, vertexCount)
	colors := make([]mgl32.Vec4, vertexCount)
	attributes := [4]meshAttribute{
		{
			location: positionLocation,
			size:     4, dim: 3, typ: gl.FLOAT,
			norm: false, data: gl.Ptr(positions),
		},
		{
			location: normalLocation,
			size:     4, dim: 3, typ: gl.FLOAT,
			norm: true, data: gl.Ptr(normals),
		},
		{
			location: texcoordsLocation,
			size:     4, dim: 2, typ: gl.FLOAT,
			norm: false, data: gl.Ptr(texcoords),
		},
		{
			location: colorLocation,
			size:     4, dim: 4, typ: gl.FLOAT,
			norm: false, data: gl.Ptr(colors),
		},
	}
	if idxs == nil {
		for i, vx := range vxs {
			positions[i] = vec3Trunc(vx.Position)
			normals[i] = vec3Trunc(vx.Normal)
			texcoords[i] = vec2Trunc(vx.TexCoords)
			colors[i] = vec4Trunc(vx.Color)
		}
	} else {
		for i, idx := range idxs {
			vx := vxs[idx]
			positions[i] = vec3Trunc(vx.Position)
			normals[i] = vec3Trunc(vx.Normal)
			texcoords[i] = vec2Trunc(vx.TexCoords)
			colors[i] = vec4Trunc(vx.Color)
		}
	}
	api.Invoke(func(api graal.Api) {
		gl.BindVertexArray(vao)
		defer gl.BindVertexArray(0)
		for i, attr := range attributes {
			size := int(attr.size) * int(attr.dim) * int(vertexCount)
			gl.EnableVertexAttribArray(attr.location)
			gl.BindBuffer(gl.ARRAY_BUFFER, vba[i])
			gl.BufferData(gl.ARRAY_BUFFER, size, attr.data, gl.STATIC_DRAW)
			gl.VertexAttribPointer(attr.location, attr.dim, attr.typ, attr.norm, 0, nil)
		}
	})
}

func glDrawVertexArray(api graal.Api, vao uint32, mode uint32, from int32, to int32) {
	if vao != 0 && to > from {
		api.Invoke(func(api graal.Api) {
			gl.BindVertexArray(vao)
			defer gl.BindVertexArray(0)
			gl.DrawArrays(mode, from, to)
		})
	}
}

func glDeleteVertexArray(api graal.Api, vao uint32) {
	if vao != 0 {
		api.Invoke(func(api graal.Api) {
			gl.DeleteVertexArrays(1, &vao)
		})
	}
}

func glDeleteBufferArray(api graal.Api, vba [4]uint32) {
	if len(vba) > 0 && vba[0] != 0 {
		api.Invoke(func(api graal.Api) {
			gl.DeleteBuffers(int32(len(vba)), &vba[0])
		})
	}
}
