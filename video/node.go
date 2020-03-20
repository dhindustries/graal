package video

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/dhindustries/graal/components"

	"github.com/dhindustries/graal/memory"

	"github.com/dhindustries/graal"
)

type node struct {
	graal.Handle
	components.Transformation
	c   []interface{}
	p   graal.Node
	l   sync.RWMutex
	api *graal.Api
}

func newNode(api *graal.Api) (graal.Node, error) {
	n := &node{
		Handle: api.NewHandle(api),
		c:      make([]interface{}, 0),
		api:    api,
	}
	n.SetScale(mgl32.Vec3{1, 1, 1})
	api.Handle(api, n)
	return n, nil
}

func (n *node) Dispose() {
	memory.Dispose(n.Handle)
	n.l.Lock()
	defer n.l.Unlock()
	if n.c != nil {
		for _, v := range n.c {
			n.releaseChild(v)
		}
	}
	n.c = nil
	n.p = nil
}

func (n *node) ParentNode() graal.Node {
	n.l.RLock()
	defer n.l.RUnlock()
	return n.p
}

func (n *node) Attach(v interface{}) {
	n.l.Lock()
	defer n.l.Unlock()
	if w, ok := v.(*node); ok {
		w.p = n
	} else if n.api.SetParent != nil {
		n.api.SetParent(n.api, v, n)
	}
	graal.Acquire(v)
	n.c = append(n.c, v)
}

func (n *node) Detach(v interface{}) {
	n.l.Lock()
	defer n.l.Unlock()
	var idx int
	var found bool
	for i, w := range n.c {
		if w == v {
			idx = i
			found = true
			break
		}
	}
	if found {
		c := n.c
		n.c = append(c[:idx], c[idx+1:])
		n.releaseChild(c)
	}
}

func (n *node) releaseChild(v interface{}) {
	if w, ok := v.(*node); ok {
		w.p = nil
	} else if n.api.SetParent != nil {
		n.api.SetParent(n.api, v, nil)
	}
	graal.Release(v)
}

func (n *node) List() []interface{} {
	n.l.RLock()
	defer n.l.RUnlock()
	return n.c
}
