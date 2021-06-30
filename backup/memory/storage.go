package memory

import "sync"

type Storage struct {
	v []Reference
	l sync.Mutex
}

func (storage *Storage) Put(v Reference) {
	storage.l.Lock()
	defer storage.l.Unlock()
	if storage.v == nil {
		storage.v = make([]Reference, 0, 1)
	}
	storage.v = append(storage.v, v)
}

func (storage *Storage) PutMany(v []Reference) {
	storage.l.Lock()
	defer storage.l.Unlock()
	if storage.v == nil {
		storage.v = make([]Reference, 0, 1)
	}
	storage.v = append(storage.v, v...)
}

func (storage *Storage) Cleanup() {
	storage.l.Lock()
	defer storage.l.Unlock()
	if storage.v != nil {
		storage.iterate()
	}
}

func (storage *Storage) iterate() {
	var a []Reference
	w := storage.v
	for {
		a = make([]Reference, 0, len(w))
		for _, v := range w {
			if v.IsValid() {
				a = append(a, v)
			} else if d, ok := v.(Disposer); ok {
				d.Dispose()
			}
		}
		if len(w) == len(a) {
			break
		}
		w = a
	}
	storage.v = a
}
