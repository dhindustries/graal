package core

import (
	"github.com/dhindustries/graal"
	mem "github.com/dhindustries/graal/memory"
)

type memory struct {
	s mem.Storage
}

func (memory *memory) handle(api *graal.Api, h graal.Handle) {
	memory.s.Put(h)
}

func (memory *memory) cleanup(api *graal.Api) {
	memory.s.Cleanup()
}
