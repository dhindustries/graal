package opengl

import (
	"sync"
	"unsafe"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"

	"github.com/dhindustries/graal/memory"
	"github.com/dhindustries/graal/utils"

	"github.com/dhindustries/graal"
	"github.com/go-gl/gl/v4.6-core/gl"
)

const (
	positionLocation  = 0
	normalLocation    = 1
	texcoordsLocation = 2
	colorLocation     = 3
)

type mesh struct {
	graal.Handle
	vao    uint32
	vba    [4]uint32
	mode   uint32
	points []graal.Vertex
	len    uint32
	valid  bool
	api    *graal.Api
	l      sync.Mutex
}

type meshResource struct {
	graal.Resource
	mesh
}

type meshAttributes struct {
	location  uint32
	size, dim uint
	ty        uint32
	norm      bool
	data      unsafe.Pointer
}

type meshRenderCommand struct {
	graal.Mesh
	graal.Texture
	graal.Color
	Transform mgl64.Mat4
}

func (m *mesh) SetVertexes(v []graal.Vertex) {
	m.l.Lock()
	defer m.l.Unlock()
	m.points = v
	m.valid = false
	m.len = uint32(len(v))
}

func (m *mesh) Build() error {
	m.l.Lock()
	defer m.l.Unlock()
	if !m.valid {
		buildMesh(m.api, m.vao, m.vba, m.points)
		m.valid = true
	}
	return nil
}

func (m *mesh) Dispose() {
	memory.Dispose(m.Handle)
	m.l.Lock()
	defer m.l.Unlock()
	vao := m.vao
	vba := m.vba
	m.vao = 0
	m.vba = [4]uint32{}
	deleteVba(m.api, vba)
	deleteVao(m.api, vao)
	m.points = nil
	m.len = 0
	m.valid = true
}

func newMesh(api *graal.Api, vx []graal.Vertex) (graal.Mesh, error) {
	ms := &mesh{
		Handle: api.NewHandle(api),
		vao:    createVao(api),
		vba:    createVba(api),
		mode:   gl.TRIANGLES,
		api:    api,
	}
	if vx != nil {
		ms.len = uint32(len(vx))
		buildMesh(ms.api, ms.vao, ms.vba, vx)
		ms.valid = true
	}
	api.Handle(api, ms)
	return ms, nil
}

func renderMesh(api *graal.Api, mi graal.Mesh) {
	var m *mesh
	switch v := mi.(type) {
	case *mesh:
		m = v
	case *meshResource:
		m = &v.mesh
	}
	if m != nil {
		m.l.Lock()
		if !m.valid {
			buildMesh(api, m.vao, m.vba, m.points)
			m.valid = true
		}
		gl.BindVertexArray(m.vao)
		defer gl.BindVertexArray(0)
		gl.DrawArrays(m.mode, 0, int32(m.len))
		m.l.Unlock()
	}
}

func createVao(api *graal.Api) uint32 {
	res := make(chan uint32)
	defer close(res)
	api.Schedule(func() {
		var vao uint32
		gl.CreateVertexArrays(1, &vao)
		res <- vao
	})
	return <-res
}

func createVba(api *graal.Api) [4]uint32 {
	res := make(chan [4]uint32)
	defer close(res)
	api.Schedule(func() {
		vba := [4]uint32{}
		gl.CreateBuffers(4, &vba[0])
		res <- vba
	})
	return <-res
}

func buildMesh(api *graal.Api, vao uint32, vba [4]uint32, vx []graal.Vertex) {
	if vx == nil {
		return
	}
	pd := make([]mgl32.Vec3, len(vx))
	nd := make([]mgl32.Vec3, len(vx))
	td := make([]mgl32.Vec2, len(vx))
	cd := make([]mgl32.Vec4, len(vx))
	attrs := [4]meshAttributes{
		meshAttributes{
			location: positionLocation,
			size:     4, dim: 3, ty: gl.FLOAT,
			norm: false, data: gl.Ptr(pd),
		},
		meshAttributes{
			location: normalLocation,
			size:     4, dim: 3, ty: gl.FLOAT,
			norm: true, data: gl.Ptr(nd),
		},
		meshAttributes{
			location: texcoordsLocation,
			size:     4, dim: 2, ty: gl.FLOAT,
			norm: false, data: gl.Ptr(td),
		},
		meshAttributes{
			location: colorLocation,
			size:     4, dim: 4, ty: gl.FLOAT,
			norm: true, data: gl.Ptr(cd),
		},
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i, v := range vx {
			pd[i] = utils.TruncVec3(v.Position)
			nd[i] = utils.TruncVec3(v.Normal)
			td[i] = utils.TruncVec2(v.TexCoords)
			cd[i] = utils.TruncVec4(v.Color)
		}
		wg.Done()
	}()
	api.Invoke(func() {
		wg.Wait()
		gl.BindVertexArray(vao)
		defer gl.BindVertexArray(0)
		for i, attr := range attrs {
			l := attr.size * attr.dim * uint(len(vx))
			gl.EnableVertexAttribArray(attr.location)
			gl.BindBuffer(gl.ARRAY_BUFFER, vba[i])
			gl.BufferData(gl.ARRAY_BUFFER, int(l), attr.data, gl.STATIC_DRAW)
			gl.VertexAttribPointer(attr.location, int32(attr.dim), attr.ty, false, 0, nil)
		}
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	})
}

func deleteVao(api *graal.Api, vao uint32) {
	if vao != 0 {
		api.Invoke(func() {
			gl.DeleteVertexArrays(1, &vao)
		})
	}
}

func deleteVba(api *graal.Api, vba [4]uint32) {
	if len(vba) > 0 && vba[0] != 0 {
		api.Invoke(func() {
			gl.DeleteBuffers(int32(len(vba)), &vba[0])
		})
	}
}
