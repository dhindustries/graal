package graal

import "time"

type Clock struct {
	current, previous time.Time
}

func (clock *Clock) Reset() {
	clock.current = time.Now()
	clock.previous = clock.current
}

func (clock *Clock) Elapsed() time.Duration {
	dur := clock.current.Sub(clock.previous)
	clock.previous = clock.current
	clock.current = time.Now()
	return dur
}
