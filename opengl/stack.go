package opengl

import (
	"sync"
)

type stack struct {
	d []interface{}
	l sync.RWMutex
}

func (s *stack) Push(v interface{}) {
	s.l.Lock()
	defer s.l.Unlock()
	if s.d == nil {
		s.d = make([]interface{}, 0, 1)
	}
	s.d = append(s.d, v)
}

func (s *stack) Top() interface{} {
	s.l.RLock()
	defer s.l.RUnlock()
	if s.d != nil && len(s.d) > 0 {
		return s.d[len(s.d)-1]
	}
	return nil
}

func (s *stack) Empty() bool {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.d != nil && len(s.d) != 0
}

func (s *stack) Pop() interface{} {
	return s.Popn(1)
}

func (s *stack) Popn(n int) interface{} {
	if n <= 0 {
		return s.Top()
	}
	s.l.Lock()
	defer s.l.Unlock()
	l := len(s.d)
	if s.d == nil || l == 0 {
		return nil
	}
	if n >= l {
		v := s.d[l-1]
		s.d = make([]interface{}, 0)
		return v
	}
	v := s.d[l-n]
	s.d = s.d[:n]
	return v
}
