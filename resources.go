package graal

import (
	"fmt"
	"log"
	"sync"
)

type Resource interface {
	Handle
	Name() string
	Path() string
	Mime() Mime
}

type ResourceManager interface {
	Initialize(engine *Engine) error
	Cleanup()
	Load(mime Mime, name string) (Resource, error)
	Register(mime Mime, loader ResourceLoader) error
	Dispose()
}

type Resources struct {
	Manager ResourceManager
}

type ResourceLoader func(resource Resource, manager ResourceManager) (Resource, error)

type resource struct {
	BaseHandle
	name, path string
	mime       Mime
}

type resourceReference struct {
	resource Resource
	err      error
	ready    bool
	lock     sync.Mutex
	cond     *sync.Cond
}

type BaseResourceManager struct {
	loaders                  map[Mime]ResourceLoader
	resources                map[Mime]map[string]*resourceReference
	logger                   *log.Logger
	loaderLock, resourceLock sync.RWMutex
}

func (resource *resource) Path() string {
	return resource.path
}

func (resource *resource) Name() string {
	return resource.name
}

func (resource *resource) Mime() Mime {
	return resource.mime
}

func (resources Resources) Load(mime Mime, path string) (Resource, error) {
	return resources.Manager.Load(mime, path)
}

func (manager *BaseResourceManager) Initialize(engine *Engine) error {
	manager.loaders = make(map[Mime]ResourceLoader)
	manager.resources = make(map[Mime]map[string]*resourceReference)
	manager.logger = engine.Logger
	manager.Register(MimeFile, fileLoader)
	manager.Register(MimeImage, imageLoader)
	return nil
}

func (manager *BaseResourceManager) Register(mime Mime, loader ResourceLoader) error {
	manager.loaderLock.Lock()
	defer manager.loaderLock.Unlock()
	mime = mime.SplitParams()
	if _, exists := manager.loaders[mime]; exists {
		return fmt.Errorf("Resource loader for %s is already registered", mime)
	}
	manager.log(fmt.Sprintf("Registered %s resource loader", mime))
	manager.loaders[mime] = loader
	return nil
}

func (manager *BaseResourceManager) Load(mime Mime, name string) (Resource, error) {
	ref, exists := manager.reference(mime, name)
	if !exists {
		if err := manager.loadInto(ref, mime, name); err != nil {
			return nil, err
		}
	}
	res, err := manager.resource(ref)
	if err != nil {
		res = nil
	}
	if res != nil {
		err = nil
		res.Acquire()
	}
	return res, err
}

func (manager *BaseResourceManager) Cleanup() {
	manager.resourceLock.Lock()
	defer manager.resourceLock.Unlock()
	manager.log("Resources cleanup...")
	for mime, list := range manager.resources {
		for name, ref := range list {
			if !manager.active(ref) {
				manager.log(fmt.Sprintf("Cleaning [%s]: %s", mime, name))
				manager.release(ref)
				delete(list, name)
			}
		}
	}
}

func (manager *BaseResourceManager) Dispose() {
	manager.resourceLock.Lock()
	manager.loaderLock.Lock()
	defer manager.resourceLock.Unlock()
	defer manager.loaderLock.Unlock()
	for _, list := range manager.resources {
		for _, ref := range list {
			manager.release(ref)
		}
	}
	manager.loaders = nil
	manager.resources = nil
	manager.logger = nil
}

func (manager *BaseResourceManager) reference(mime Mime, name string) (*resourceReference, bool) {
	manager.resourceLock.Lock()
	defer manager.resourceLock.Unlock()
	if _, exists := manager.resources[mime]; !exists {
		manager.resources[mime] = make(map[string]*resourceReference)
	}
	ref, exists := manager.resources[mime][name]
	if !exists {
		ref = &resourceReference{}
		ref.cond = sync.NewCond(&ref.lock)
		manager.resources[mime][name] = ref
	}
	return ref, exists
}

func (manager *BaseResourceManager) loadInto(ref *resourceReference, mime Mime, name string) error {
	loader, err := manager.loader(mime)
	if err != nil {
		return err
	}
	manager.log(fmt.Sprintf("Loading resource [%s]: %s", mime, name))
	go func() {
		base := &resource{}
		base.name = name
		base.path = name
		base.mime = mime
		res, err := loader(base, manager)
		ref.lock.Lock()
		ref.resource = res
		ref.err = err
		ref.ready = true
		ref.lock.Unlock()
		ref.cond.Broadcast()
	}()
	return nil
}

func (manager *BaseResourceManager) resource(ref *resourceReference) (Resource, error) {
	ref.lock.Lock()
	defer ref.lock.Unlock()
	if ref.ready {
		return ref.resource, ref.err
	}
	ref.cond.Wait()
	return ref.resource, ref.err
}

func (manager *BaseResourceManager) active(ref *resourceReference) bool {
	ref.lock.Lock()
	defer ref.lock.Unlock()
	return !ref.ready || (ref.resource != nil && ref.resource.Valid())
}

func (manager *BaseResourceManager) release(ref *resourceReference) {
	if ref.resource != nil {
		if res, ok := ref.resource.(Disposable); ok {
			res.Dispose()
		}
	}
}

func (manager *BaseResourceManager) loader(mime Mime) (ResourceLoader, error) {
	manager.loaderLock.RLock()
	defer manager.loaderLock.RUnlock()

	mime = mime.SplitParams()
	mimes := []Mime{
		mime,
		mime.SplitSubType(),
	}
	for _, mime := range mimes {
		if loader, exists := manager.loaders[mime]; exists {
			return loader, nil
		}
	}
	return nil, fmt.Errorf("Resource loader for %s not registered", mime)
}

func (manager *BaseResourceManager) log(v interface{}) {
	if manager.logger != nil {
		manager.logger.Println(v)
	}
}
