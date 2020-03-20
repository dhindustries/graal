package action

import (
	"sync"
	"time"
)

type Scheduler struct {
	seq   *Sequence
	tick  *time.Ticker
	close chan bool
	l     sync.Mutex
}

func NewScheduler() *Scheduler {
	return &Scheduler{seq: NewSequence()}
}

func (s *Scheduler) Add(action Action) {
	s.seq.Add(action)
}

func (s *Scheduler) Start(t interface{}, freq float32) {
	s.l.Lock()
	if s.close != nil {
		s.l.Unlock()
		return
	}
	s.close = make(chan bool)
	s.l.Unlock()
	s.tick = time.NewTicker(time.Second / time.Duration(freq))
	pt := time.Now()
loop:
	for {
		select {
		case <-s.close:
			s.tick.Stop()
			break loop
		case ct := <-s.tick.C:
			s.Update(t, ct.Sub(pt))
			pt = ct
		}
	}
	s.l.Lock()
	defer s.l.Unlock()
	s.tick.Stop()
	close(s.close)
	s.close = nil
	s.tick = nil
}

func (s *Scheduler) Stop() {
	s.l.Lock()
	defer s.l.Unlock()
	if s.close != nil {
		s.close <- true
	}
}

func (s *Scheduler) Dispose() {
	s.Stop()
}

func (s *Scheduler) Update(t interface{}, dt time.Duration) {
	if s.seq.Run(t, dt) {
		s.Stop()
	}
}
