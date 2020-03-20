package core

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/dhindustries/graal"
)

var parents map[uintptr]*sync.Pool
var parentsLock sync.Mutex

func objID(data interface{}) uintptr {
	if data == nil {
		return uintptr(unsafe.Pointer(nil))
	}
	var addr uintptr
	v := reflect.ValueOf(data)
	switch v.Type().Kind() {
	case reflect.Ptr:
		addr = v.Elem().UnsafeAddr()
	case reflect.Uintptr:
		addr = v.Pointer()
	}
	return addr
}

func setParent(api *graal.Api, v interface{}, p graal.Node) {
	id := objID(v)
	parentsLock.Lock()
	defer parentsLock.Unlock()
	if parents == nil {
		parents = make(map[uintptr]*sync.Pool)
	}
	var pool *sync.Pool
	if pool, ok := parents[id]; !ok {
		pool = new(sync.Pool)
		parents[id] = pool
	}
	pool.Put(v)
}

func getParent(api *graal.Api, v interface{}) graal.Node {
	id := objID(v)
	parentsLock.Lock()
	defer parentsLock.Unlock()
	if parents != nil {
		if pool, ok := parents[id]; ok {
			n := pool.Get().(graal.Node)
			pool.Put(n)
			return n
		}
	}
	return nil
}
