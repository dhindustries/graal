package opengl

import (
	"github.com/dhindustries/graal"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type renderer struct {
	prog graal.Program
	cam  graal.Camera
	wnd  graal.Window
	q    queue
}

type renderItem struct {
	V interface{}
	T mgl32.Mat4
}

type swapper interface {
	Swap()
}

type textured interface {
	Texture() graal.Texture
}

type meshed interface {
	Mesh() graal.Mesh
}

type shaped interface {
	Shape() graal.Shape
}

type transformed interface {
	Transform() mgl32.Mat4
}

func (r *renderer) commit(org *graal.Api) {
	org.Invoke(func() {
		api := graal.Fork(org)
		api.Invoke = func(fn func()) {
			fn()
		}
		r.render(api)
		if w, ok := r.wnd.(swapper); ok {
			w.Swap()
		}
	})
}

func (r *renderer) render(api *graal.Api) {
	r.q.l.Lock()
	defer r.q.l.Unlock()
	r.begin(api)
	if r.q.d != nil {
		for _, i := range r.q.d {
			switch cmd := i.(type) {
			case meshRenderCommand:
				r.setModel(api, cmd.Transform)
				bindTexture(api, cmd.Texture)
				renderMesh(api, cmd.Mesh)
			}
		}
		r.q.d = make([]interface{}, 0)
	}
	r.end(api)
}

func (r *renderer) enqueue(api *graal.Api, elem interface{}, t mgl32.Mat4) {
	if node, ok := elem.(graal.Node); ok {
		tf := t.Mul4(node.Transform())
		for _, child := range node.List() {
			api.RenderEnqueue(api, child, tf)
		}
	} else if obj, ok := elem.(transformed); ok {
		t = t.Mul4(obj.Transform())
	}
	var tex graal.Texture
	if v, ok := elem.(textured); ok {
		tex = v.Texture()
	}
	if v, ok := elem.(meshed); ok {
		elem = v.Mesh()
	}
	switch v := elem.(type) {
	case graal.Mesh:
		r.q.Push(meshRenderCommand{
			Mesh:      v,
			Texture:   tex,
			Transform: t,
		})
	}
}

func (r *renderer) finit(api *graal.Api) {
	r.bindProgram(api, nil)
	r.setCamera(api, nil)
}

func (r *renderer) begin(api *graal.Api) {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	var pid uint32
	switch v := r.prog.(type) {
	case *program:
		pid = v.id
	case *programResource:
		pid = v.program.id
	}
	gl.UseProgram(pid)
	if pid != 0 {
		if r.cam != nil {
			r.setView(api, r.cam.View())
			r.setProjection(api, r.cam.Projection())
		} else {
			r.setView(api, mgl32.Ident4())
			r.setProjection(api, mgl32.Ident4())
		}
	}
}

func (r *renderer) end(api *graal.Api) {

}

func (r *renderer) bindProgram(api *graal.Api, prog graal.Program) {
	if r.prog != nil {
		r.prog.Release()
	}
	if prog != nil {
		prog.Acquire()
	}
	r.prog = prog
}

func (r *renderer) setCamera(api *graal.Api, cam graal.Camera) {
	if r.cam != nil {
		r.cam.Release()
	}
	if cam != nil {
		cam.Acquire()
	}
	r.cam = cam
}
