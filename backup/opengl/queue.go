package opengl

import (
	"sync"
)

type queue struct {
	d []interface{}
	l sync.RWMutex
}

func (q *queue) Push(v interface{}) {
	q.l.Lock()
	defer q.l.Unlock()
	if q.d == nil {
		q.d = make([]interface{}, 0, 1)
	}
	q.d = append(q.d, v)
}

func (q *queue) PushFront(v interface{}) {
	q.l.Lock()
	defer q.l.Unlock()
	if q.d == nil {
		q.d = make([]interface{}, 0, 1)
	}
	q.d = append([]interface{}{v}, q.d)
}

func (q *queue) Top() interface{} {
	q.l.RLock()
	defer q.l.RUnlock()
	if q.d != nil && len(q.d) > 0 {
		return q.d[0]
	}
	return nil
}

func (q *queue) Empty() bool {
	q.l.RLock()
	defer q.l.RUnlock()
	return q.d != nil && len(q.d) != 0
}

func (q *queue) Pop() interface{} {
	return q.Popn(1)
}

func (q *queue) Popn(n int) interface{} {
	if n <= 0 {
		return q.Top()
	}
	q.l.Lock()
	defer q.l.Unlock()
	l := len(q.d)
	if q.d == nil || l == 0 {
		return nil
	}
	if n >= l {
		v := q.d[l-1]
		q.d = make([]interface{}, 0)
		return v
	}
	v := q.d[n]
	q.d = q.d[n:]
	return v
}
