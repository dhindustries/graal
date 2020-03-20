package core

import (
	"fmt"
	"sync"

	"github.com/dhindustries/graal"
)

type resource struct {
	graal.Handle
	name, path string
	mime       graal.Mime
	logfn      func(string, ...interface{})
}

type resourceLoader func(*graal.Api, graal.Resource) (graal.Resource, error)
type resourceReference struct {
	mime  graal.Mime
	name  string
	r     graal.Resource
	e     error
	l     sync.Mutex
	c     *sync.Cond
	logfn func(string, ...interface{})
}

type resourceManager struct {
	lo     map[graal.Mime]resourceLoader
	rs     map[graal.Mime]map[string]*resourceReference
	ll, rl sync.RWMutex
}

func (r *resource) Path() string {
	return r.path
}

func (r *resource) Name() string {
	return r.name
}

func (r *resource) Mime() graal.Mime {
	return r.mime
}

func (r *resource) Dispose() {
	if r.logfn != nil {
		r.logfn("Disposing resource [%s] %s\n", r.mime, r.name)
	}
}

func (ref *resourceReference) get() (graal.Resource, error) {
	ref.l.Lock()
	defer ref.l.Unlock()
	if ref.r != nil || ref.e != nil {
		if ref.r != nil {
			ref.r.Acquire()
		}
		return ref.r, ref.e
	}
	ref.c.Wait()
	if ref.r != nil {
		ref.r.Acquire()
	}
	return ref.r, ref.e
}

func (ref *resourceReference) set(r graal.Resource, e error) {
	ref.l.Lock()
	ref.r = r
	ref.e = e
	if ref.logfn != nil {
		if e != nil {
			ref.logfn("Failed to load resource [%s] %s\n", ref.mime, ref.name)
		} else if r != nil {
			ref.logfn("Loaded resource [%s] %s\n", ref.mime, ref.name)
		}
	}
	ref.l.Unlock()
	ref.c.Broadcast()
}

func (ref *resourceReference) valid() bool {
	ref.l.Lock()
	defer ref.l.Unlock()
	return (ref.r == nil && ref.e == nil) || (ref.r != nil && ref.r.IsValid())
}

func newResourceManager() *resourceManager {
	return &resourceManager{
		lo: make(map[graal.Mime]resourceLoader),
		rs: make(map[graal.Mime]map[string]*resourceReference),
	}
}

func (manager *resourceManager) loadResource(api *graal.Api, mime graal.Mime, path string) (graal.Resource, error) {
	ref, ex := manager.reference(api, mime, path)
	if !ex {
		l, err := manager.loader(mime)
		if err != nil {
			return nil, err
		}
		r, err := manager.newResource(api, mime, path)
		if err != nil {
			return nil, err
		}
		manager.log("Loading resource [%s] %s\n", mime, path)
		go func() {
			v, err := l(api, r)
			if v != nil {
				api.Handle(api, v)
			}
			ref.set(v, err)
		}()
	}
	return ref.get()
}

func (manager *resourceManager) reference(api *graal.Api, mime graal.Mime, name string) (*resourceReference, bool) {
	manager.rl.Lock()
	defer manager.rl.Unlock()
	if manager.rs == nil {
		manager.rs = make(map[graal.Mime]map[string]*resourceReference)
	}
	if _, ok := manager.rs[mime]; !ok {
		manager.rs[mime] = make(map[string]*resourceReference)
	}
	v, ok := manager.rs[mime][name]
	if !ok {
		v = &resourceReference{}
		v.mime = mime
		v.name = name
		v.c = sync.NewCond(&v.l)
		v.logfn = manager.log
		manager.rs[mime][name] = v
	}
	return v, ok
}

func (manager *resourceManager) newResource(api *graal.Api, mime graal.Mime, path string) (graal.Resource, error) {
	return &resource{
		Handle: &handle{},
		name:   path,
		path:   path,
		mime:   mime,
		logfn:  manager.log,
	}, nil
}

func (manager *resourceManager) cleanup(api *graal.Api) {
	manager.rl.Lock()
	defer manager.rl.Unlock()
	manager.log("Resources cleanup...\n")

	pass := func(d map[graal.Mime]map[string]*resourceReference, c *int) map[graal.Mime]map[string]*resourceReference {
		for m, l := range d {
			for p, r := range l {
				if r.valid() {
					*c++
				} else {
					delete(d[m], p)
				}
			}
		}
		return d
	}
	cl := 0
	d := pass(manager.rs, &cl)
	for {
		cc := 0
		d = pass(d, &cc)
		if cc == cl {
			break
		}
		cl = cc
	}
	manager.rs = d
}

func (manager *resourceManager) setLoader(api *graal.Api, mime graal.Mime, loader func(*graal.Api, graal.Resource) (graal.Resource, error)) {
	manager.ll.Lock()
	defer manager.ll.Unlock()
	manager.lo[mime] = loader
}

func (manager *resourceManager) loader(mime graal.Mime) (resourceLoader, error) {
	manager.ll.RLock()
	defer manager.ll.RUnlock()
	mime = mime.SplitParams()
	mimes := []graal.Mime{
		mime,
		mime.SplitSubType(),
	}
	for _, mime := range mimes {
		if loader, exists := manager.lo[mime]; exists {
			return loader, nil
		}
	}
	return nil, fmt.Errorf("Resource loader for %s not registered", mime)
}

func (manager *resourceManager) log(f string, v ...interface{}) {
	fmt.Printf(f, v...)
}
